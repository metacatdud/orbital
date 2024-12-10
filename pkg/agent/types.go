package agent

type Network struct {
	Name string
	Type string
}

type Cmd struct {
	Name  string
	Value string
}

type Port struct {
	Name     string
	Internal string
	External string
}

type Volume struct {
	Name          string
	ContainerPath string
	HostPath      string
}

type EnvVar struct {
	Name  string
	Value string
}

type Label struct {
	Name  string
	Value string
}
