package agent

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"orbital/pkg/logger"
	"strings"
)

type Docker struct {
	client *client.Client
	log    *logger.Logger
}

func (agent *Docker) CreateContainer(containerCfg Container) (string, error) {

	// Check network
	if err := agent.ensureNetworkExist(containerCfg.Network); err != nil {
		return "", err
	}

	// Ports
	portBindings, exposedPorts, err := convertPorts(containerCfg.Ports)
	if err != nil {
		return "", err
	}

	// Labels
	labels := map[string]string{
		"orbital.managed": "true",
	}

	for _, label := range containerCfg.Labels {
		labels[label.Name] = label.Value
	}

	// Environment vars
	var env []string
	for _, envVar := range containerCfg.EnvVars {
		env = append(env, fmt.Sprintf("%s=%s", envVar.Name, envVar.Value))
	}

	// Volumes
	volumes := convertVolumes(containerCfg.Volumes)

	// Container configuration
	config := &container.Config{
		Image:        containerCfg.Image,
		Cmd:          convertCmdsToString(containerCfg.Cmds),
		ExposedPorts: exposedPorts,
		Labels:       labels,
		Env:          env,
	}

	// Host configuration
	hostConfig := &container.HostConfig{
		AutoRemove:   true,
		PortBindings: portBindings,
		Binds:        volumes,
	}

	// Network configuration
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			containerCfg.Network.Name: {},
		},
	}

	res, err := agent.client.ContainerCreate(context.Background(), config, hostConfig, networkConfig, nil, containerCfg.Name)
	if err != nil {
		return "", err
	}

	agent.log.Info("Created container. NOT Running", "id", res.ID)

	return res.ID, nil
}

func (agent *Docker) ListContainers() ([]*Container, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "orbital.managed=true")

	opts := container.ListOptions{
		All:     true,
		Filters: filterArgs,
	}

	containers, err := agent.client.ContainerList(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	var orbitalContainers []*Container
	for _, c := range containers {
		var ports []Port
		for _, p := range c.Ports {
			ports = append(ports, Port{
				Name:     "", // No name is available in this context. Update later form db if needed
				Internal: fmt.Sprintf("%d", p.PrivatePort),
				External: fmt.Sprintf("%d", p.PublicPort),
			})
		}

		var labels []Label
		for key, value := range c.Labels {
			labels = append(labels, Label{Name: key, Value: value})
		}

		orbitalContainers = append(orbitalContainers, &Container{
			Name:   c.Names[0],
			Image:  c.Image,
			Ports:  ports,
			Labels: labels,
		})
	}

	return orbitalContainers, nil
}

func (agent *Docker) GetContainer(id string) (*Container, error) {
	c, err := agent.client.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}
	// Map metadata to `Container` struct
	var ports []Port
	for port, bindings := range c.NetworkSettings.Ports {
		for _, binding := range bindings {
			ports = append(ports, Port{
				Name:     "", // No name is available in this context
				Internal: port.Port(),
				External: binding.HostPort,
			})
		}
	}

	var labels []Label
	for key, value := range c.Config.Labels {
		labels = append(labels, Label{Name: key, Value: value})
	}

	var envVars []EnvVar
	for _, env := range c.Config.Env {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envVars = append(envVars, EnvVar{Name: parts[0], Value: parts[1]})
		}
	}

	var volumes []Volume
	for _, mount := range c.Mounts {
		volumes = append(volumes, Volume{
			Name:          mount.Name,
			HostPath:      mount.Source,
			ContainerPath: mount.Destination,
		})
	}

	return &Container{
		Image: c.Config.Image,
		Name:  c.Name[1:],
		Network: Network{
			Name: c.NetworkSettings.Networks["bridge"].NetworkID,
			Type: "bridge",
		},
		Cmds:    convertCmdsFromString(c.Config.Cmd),
		Ports:   ports,
		Volumes: volumes,
		EnvVars: envVars,
		Labels:  labels,
	}, nil
}

func (agent *Docker) ensureNetworkExist(net Network) error {
	opts := network.ListOptions{}

	networks, err := agent.client.NetworkList(context.Background(), opts)
	if err != nil {
		return err
	}

	for _, nw := range networks {
		if nw.Name == net.Name {
			return nil
		}
	}

	_, err = agent.client.NetworkCreate(context.Background(), net.Name, network.CreateOptions{
		Driver:     network.NetworkBridge,
		Scope:      "local",
		EnableIPv6: PrimitiveToPtr(true),
		Attachable: true,
		Labels: map[string]string{
			"orbital.managed": "true",
		},
		IPAM: &network.IPAM{
			Driver: "default",
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func NewDocker() (*Docker, error) {
	lg := logger.New(logger.LevelDebug, logger.FormatString)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Docker{
		client: cli,
		log:    lg,
	}, nil
}

func PrimitiveToPtr[T any](val T) *T {
	return &val
}

func convertPorts(ports []Port) (nat.PortMap, nat.PortSet, error) {
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}

	for _, port := range ports {
		natPort, err := nat.NewPort("tcp", port.Internal)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid container port %s: %w", port.Internal, err)
		}

		// Add to PortBindings
		portBindings[natPort] = []nat.PortBinding{
			{HostPort: port.External},
		}

		// Add to ExposedPorts
		exposedPorts[natPort] = struct{}{}
	}

	return portBindings, exposedPorts, nil
}

func convertVolumes(volumes []Volume) []string {
	var binds []string
	for _, volume := range volumes {
		binds = append(binds, fmt.Sprintf("%s:%s", volume.HostPath, volume.ContainerPath))
	}
	return binds
}

func convertCmdsToString(cmds []Cmd) []string {
	var convertedCmds []string
	for _, cmd := range cmds {
		convertedCmds = append(convertedCmds, cmd.Name)
	}
	return convertedCmds
}

func convertCmdsFromString(cmds []string) []Cmd {
	var extractedCmds []Cmd
	for _, cmd := range cmds {
		extractedCmds = append(extractedCmds, Cmd{Name: cmd})
	}
	return extractedCmds
}
