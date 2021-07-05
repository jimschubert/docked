package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

func maintainerDeprecated() validations.Rule {
	r := validations.SimpleRegexRule{
		Name:     "maintainer-deprecated",
		Summary:  "MAINTAINER is deprecated",
		Details:  "MAINTAINER instruction is deprecated; Use LABEL instead, which can be queried via `docker inspect`.",
		Pattern:  `[[:graph:]]+`,
		Priority: model.LowPriority,
		Command:  commands.Maintainer,
		URL:      model.StringPtr("https://docs.docker.com/engine/reference/builder/#maintainer-deprecated"),
	}
	return &r
}

func init() {
	AddRule(maintainerDeprecated())
}
