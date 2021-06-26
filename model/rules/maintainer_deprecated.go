package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

func maintainerDeprecated() validations.Rule {
	return validations.NewSimpleRegexRule(
		"maintainer-deprecated",
		"MAINTAINER instruction is deprecated; Use LABEL instead, which can be queried via 'docker inspect'.",
		`[[:graph:]]+`,
		model.LowPriority,
		commands.Maintainer,
		nil,
		model.StringPtr("https://docs.docker.com/engine/reference/builder/#maintainer-deprecated"),
	)
}

func init() {
	AddRule(maintainerDeprecated())
}
