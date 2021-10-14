package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/shell"
	"github.com/jimschubert/docked/model/validations"
	log "github.com/sirupsen/logrus"
)

func gpgWithoutBatch() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "gpg-without-batch",
		Summary:  "GPG call without --batch (or --no-tty) may error.",
		Details:  "Running GPG without --batch (or --no-tty) may cause GPG to fail opening /dev/tty, resulting in docker build failures.",
		Priority: model.MediumPriority,
		Commands: []commands.DockerCommand{commands.Run},
		URL:      model.StringPtr("https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=913614"),

		Evaluator: validations.MultiContextFullEvaluator{
			Fn: func(mcr *validations.MultiContextRule) *validations.ValidationResult {
				if mcr == nil || mcr.ContextCache == nil {
					return &validations.ValidationResult{
						Result:  model.Skipped,
						Details: mcr.GetSummary(),
					}
				}
				result := model.Success

				for _, nodeContext := range *mcr.ContextCache {
					posixCommands, err := shell.NewPosixCommandFromNode(&nodeContext.Node)
					if err != nil {
						log.Warnf("Unable to parse RUN command, skipping validation: %#v", nodeContext.Node.Location())
						result = model.Skipped
					} else {
						for _, command := range posixCommands {
							if (command.Name == "gpg" || command.Name == `\gpg`) && command.Args != nil {
								var hasBatch bool
								var hasNoTty bool
								for _, arg := range command.Args {
									hasBatch = hasBatch || strings.Contains(arg, "--batch")
									hasNoTty = hasNoTty || strings.Contains(arg, "--no-tty")
								}

								if !hasBatch && !hasNoTty {
									result = model.Failure
									nodeContext.Context.CausedFailure = true
								}
							}
						}
					}
				}

				return &validations.ValidationResult{
					Result:   result,
					Details:  mcr.GetSummary(),
					Contexts: *mcr.GetContexts(),
				}
			},
		},
	}
	return &r
}

func init() {
	AddRule(gpgWithoutBatch())
}
