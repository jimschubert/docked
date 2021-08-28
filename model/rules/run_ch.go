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
		Summary:  "Change ownership in the same layer as file operation (RUN or COPY)",
		Details:  "In AUFS, ownership defined in an earlier layer can not be overridden by a broader mask in a later layer.",
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
