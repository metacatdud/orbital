package agent

type Container struct {
	Name    string
	Image   string
	Network Network
	Ports   []Port
	Cmds    []Cmd
	Volumes []Volume
	EnvVars []EnvVar
	Labels  []Label
}
