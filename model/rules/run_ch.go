package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

// Concept taken from https://github.com/docker-library/dockerfile-validator
func layeredChownChmod() validations.Rule {
	rule := validations.SimpleDeferredRegexRule{
		Name:     "layered-ownership-change",
		Summary:  "Some storage drivers may have issues with ownership changes in different layers. Move this to an earlier layer if possible.",
		Patterns: []string{`^ch(own|mod)\b`},
		Priority: model.MediumPriority,
		Commands: []commands.DockerCommand{commands.Run},
		URL:      model.StringPtr("https://github.com/moby/moby/issues/783#issuecomment-19237045"),
	}
	return &rule
}

func init() {
	AddRule(layeredChownChmod())
}
