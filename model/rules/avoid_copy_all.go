package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func avoidCopyAll() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "avoid-copy-all",
		Summary:  "Avoid copying entire source directory into image",
		Details:  "Explicitly copying sources helps avoid accidentally persisting secrets or other files that should not be shared.",
		Priority: model.HighPriority,
		Commands: []commands.DockerCommand{commands.Copy},
		Evaluator: validations.MultiContextPerNodeEvaluator{
			Fn: func(node *parser.Node, validationContext validations.ValidationContext) model.Valid {
				if strings.HasPrefix(node.Original, "COPY . ") {
					return model.Recommendation
				}
				return model.Skipped
			},
		},
	}
	return &r
}

func init() {
	AddRule(avoidCopyAll())
}
