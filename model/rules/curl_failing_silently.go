package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/shell"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

func curlWithoutFail() validations.Rule {
	r := validations.MultiContextRule{
		Name:    "curl-without-fail",
		Summary: "Avoid using curl without the silent failing option -f/--fail",
		Details: "Invoking curl without -f/--fail may result in incorrect, missing or stale data, which is a security concern. " +
			"Ignore this rule only if you're handling server errors or verifying file contents separately.",
		Priority:         model.CriticalPriority,
		Commands:         []commands.DockerCommand{commands.Run},
		AppliesToBuilder: true,
		Category:         nil,
		URL:              model.StringPtr("https://curl.se/docs/faq.html#Why_do_I_get_downloaded_data_eve"),
		Evaluator: validations.MultiContextPerNodeEvaluator{
			Fn: func(node *parser.Node, validationContext validations.ValidationContext) model.Valid {
				posixCommands, err := shell.NewPosixCommandFromNode(node)
				if err != nil {
					log.Warnf("Unable to parse RUN command, skipping validation: %#v", node.Location())
					return model.Skipped
				}
				for _, command := range posixCommands {
					if (command.Name == "curl" || command.Name == `\curl`) && command.Args != nil {
						var hasRequiredFlag bool
						for _, arg := range command.Args {
							if hasRequiredFlag {
								break
							}

							hasRequiredFlag = hasRequiredFlag || model.NewPattern(`\b-f\b|\b--fail\b`).Matches(arg)
						}

						if !hasRequiredFlag {
							return model.Failure
						}
					}
				}
				return model.Success
			},
		},
	}
	return &r
}

func init() {
	AddRule(curlWithoutFail())
}
