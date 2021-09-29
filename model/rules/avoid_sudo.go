package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/shell"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

func avoidSudo() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "avoid-sudo",
		Summary:  "Avoid running root elevation tasks like sudo/su",
		Details:  "Non-root users should avoid having sudo access in containers, as it has unpredictable TTY and signal-forwarding behavior that can cause problems. Consider using gosu instead.",
		Priority: model.MediumPriority,
		Commands: []commands.DockerCommand{commands.Run},
		URL:      model.StringPtr("https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user"),
		Evaluator: validations.MultiContextPerNodeEvaluator{
			Fn: func(node *parser.Node, validationContext validations.ValidationContext) model.Valid {
				posixCommands, err := shell.NewPosixCommandFromNode(node)
				if err != nil {
					log.Warnf("Unable to parse RUN command, skipping validation: %#v", node.Location())
					return model.Skipped
				}
				for _, command := range posixCommands {
					if command.Name == "su" || command.Name == "sudo" || command.Name == `\su` || command.Name == `\sudo` {
						return model.Recommendation
					}
				}
				return model.Skipped
			},
		},
	}
	return &r
}

func init() {
	AddRule(avoidSudo())
}
