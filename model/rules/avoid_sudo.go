package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func avoidSudo() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "avoid-sudo",
		Summary:  "Avoid running root elevation tasks like sudo/su",
		Details:  "Non-root users should avoid having sudo access in containers. Consider using gosu instead.",
		Priority: model.MediumPriority,
		Commands: []commands.DockerCommand{commands.Run},
		URL:      model.StringPtr("https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user"),
		Evaluator: func(node *parser.Node, validationContext validations.ValidationContext) model.Valid {
			trimStart := len(node.Value) + 1 // command plus trailing space
			matchAgainst := node.Original[trimStart:]
			if model.NewPattern(`\bsu\b|\bsudo\b`).Matches(matchAgainst) {
				return model.Recommendation
			}
			return model.Skipped
		},
	}
	return &r
}

func init() {
	AddRule(avoidSudo())
}
