package commands

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// DockerCommand extends the representation of a docker command
type DockerCommand string

//goland:noinspection ALL
const (
	Add         = DockerCommand("add")
	Arg         = DockerCommand("arg")
	Cmd         = DockerCommand("cmd")
	Copy        = DockerCommand("copy")
	Entrypoint  = DockerCommand("entrypoint")
	Env         = DockerCommand("env")
	Expose      = DockerCommand("expose")
	From        = DockerCommand("from")
	Healthcheck = DockerCommand("healthcheck")
	Label       = DockerCommand("label")
	Maintainer  = DockerCommand("maintainer")
	Onbuild     = DockerCommand("onbuild")
	Run         = DockerCommand("run")
	Shell       = DockerCommand("shell")
	StopSignal  = DockerCommand("stopsignal")
	User        = DockerCommand("user")
	Volume      = DockerCommand("volume")
	Workdir     = DockerCommand("workdir")
)

// Upper returns the uppercase representation of the DockerCommand value.
func (d DockerCommand) Upper() string {
	return strings.ToUpper(string(d))
}

func (d *DockerCommand) UnmarshalYAML(value *yaml.Node) error {
	var original string
	if err := value.Decode(&original); err != nil {
		return err
	}
	switch original {
	case "add":
		*d = Add
	case "arg":
		*d = Arg
	case "cmd":
		*d = Cmd
	case "copy":
		*d = Copy
	case "entrypoint":
		*d = Entrypoint
	case "env":
		*d = Env
	case "expose":
		*d = Expose
	case "from":
		*d = From
	case "healthcheck":
		*d = Healthcheck
	case "label":
		*d = Label
	case "maintainer":
		*d = Maintainer
	case "onbuild":
		*d = Onbuild
	case "run":
		*d = Run
	case "shell":
		*d = Shell
	case "stopsignal":
		*d = StopSignal
	case "user":
		*d = User
	case "volume":
		*d = Volume
	case "workdir":
		*d = Workdir
	default:
		return fmt.Errorf("unknown DockerCommand: %s", original)
	}
	return nil
}
