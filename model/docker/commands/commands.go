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

func Of(command string) DockerCommand {
	switch strings.ToLower(command) {
	case "add":
		return Add
	case "arg":
		return Arg
	case "cmd":
		return Cmd
	case "copy":
		return Copy
	case "entrypoint":
		return Entrypoint
	case "env":
		return Env
	case "expose":
		return Expose
	case "from":
		return From
	case "healthcheck":
		return Healthcheck
	case "label":
		return Label
	case "maintainer":
		return Maintainer
	case "onbuild":
		return Onbuild
	case "run":
		return Run
	case "shell":
		return Shell
	case "stopsignal":
		return StopSignal
	case "user":
		return User
	case "volume":
		return Volume
	case "workdir":
		return Workdir
	default:
		return "" // Return an empty DockerCommand for unknown commands
	}
}

// Upper returns the uppercase representation of the DockerCommand value.
func (d DockerCommand) Upper() string {
	return strings.ToUpper(string(d))
}

// UnmarshalYAML unmarshalls from a YAML node into a DockerCommand
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
