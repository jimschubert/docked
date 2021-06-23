package commands

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
