package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

func noDebianFrontend() validations.Rule {
	return validations.NewSimpleRegexRule(
		"no-debian-frontend",
		"Avoid DEBIAN_FRONTEND, which affects derived images and docker run. Change this to an ARG.",
		`DEBIAN_FRONTEND`,
		model.CriticalPriority,
		commands.Env,
		nil,
		nil,
		)
}

func init() {
	AddRule(noDebianFrontend())
}
