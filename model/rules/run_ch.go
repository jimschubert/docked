package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

// Concept taken from https://github.com/docker-library/dockerfile-validator
func layeredChownChmod() validations.Rule {
	return validations.NewSimpleDeferredRegexRule(
		"layered-ownership-change",
		"Some storage drivers may have issues with ownership changes in different layers. Move this to an earlier layer if possible.",
		[]string { `^ch(own|mod)\b` },
		model.MediumPriority,
		[]commands.DockerCommand { commands.Run },
		false,
		nil,
		model.StringPtr("https://github.com/moby/moby/issues/783#issuecomment-19237045"),
	)
}

func init() {
	AddRule(layeredChownChmod())
}
