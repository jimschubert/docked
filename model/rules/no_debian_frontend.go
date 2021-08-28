package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

func noDebianFrontend() validations.Rule {
	r := validations.SimpleRegexRule{
		Name:     "no-debian-frontend",
		Summary:  "Convert DEBIAN_FRONTEND to an ARG.",
		Details:  "Avoid DEBIAN_FRONTEND, which affects derived images and docker run. Change this to an ARG.",
		Pattern:  `\bDEBIAN_FRONTEND\b`,
		Priority: model.CriticalPriority,
		Command:  commands.Env,
	}
	return &r
}

func init() {
	AddRule(noDebianFrontend())
}
