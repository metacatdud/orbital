package main

import (
	"fmt"
	"orbital/cmd"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// TODO: Instantiate an manager with HTTP
// TODO: Onboard: Admin Public key, storage path, http and server ports
// TODO: Instantiate a SQLite handler
// TODO: Initiate a member list node
// TODO: Get machine's available resources
func run() error {
	//dockerAgent, err := agent.NewDocker()
	//if err != nil {
	//	return fmt.Errorf("new docker agent: %w", err)
	//}
	//
	//containerCfg := agent.Container{
	//	Name:  "orbital-redis",
	//	Image: "redis:7-alpine",
	//	Network: agent.Network{
	//		Name: "orbital_net",
	//	},
	//	Ports: []agent.Port{
	//		{
	//			Name:     "6379/tcp",
	//			Internal: "6379",
	//			External: "6300",
	//		},
	//	},
	//	Volumes: []agent.Volume{
	//		{
	//			Name:          "redisVol_1",
	//			ContainerPath: "/data",
	//			HostPath:      "/home/tibi/projects/orbital/storage/redis",
	//		},
	//	},
	//}
	//
	//containerId, err := dockerAgent.CreateContainer(containerCfg)
	//if err != nil {
	//	return fmt.Errorf("create container: %w", err)
	//}
	//
	//fmt.Printf("created container:%s\n", containerId)
	//
	//containers, err := dockerAgent.ListContainers()
	//if err != nil {
	//	return fmt.Errorf("list containers: %w", err)
	//}
	//
	//fmt.Printf("containers:%+v\n", containers)

	return cmd.Execute()
}
