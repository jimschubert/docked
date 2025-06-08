package rules

import (
	"regexp"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/shell"
	"github.com/jimschubert/docked/model/validations"
	log "github.com/sirupsen/logrus"
)

func considerMultistageBuild() validations.Rule {
	var buildTools = [...]string{`\b[\\]?np[mx]\b`, `\b[\\]?mvn[w]?\b`, `\b[\\]?bazel\b`, `\b[\\]?gradle[w]?\b`, `\b[\\]?go\b`, `\b[\\]?goreleaser\b`}
	r := validations.MultiContextRule{
		Name:    "consider-multistage",
		Summary: "Consider using multi-stage builds for complex operations like building code.",
		Details: "A multi-stage build can reduce the final image size by building necessary components or" +
			" downloading large archives in a separate build context. This can help keep your final image lean.",
		Priority: model.LowPriority,
		Commands: []commands.DockerCommand{commands.Run, commands.From},
		URL:      model.StringPtr("https://docs.docker.com/develop/develop-images/multistage-build/"),
		Evaluator: validations.MultiContextFullEvaluator{
			Fn: func(mcr *validations.MultiContextRule) *validations.ValidationResult {
				var hasFailures bool
				var hasAnyBuilder bool
				validationContexts := make([]validations.ValidationContext, 0)
				for _, nodeContext := range *mcr.ContextCache {
					if nodeContext.Context.IsBuilderContext {
						hasAnyBuilder = true
					}
					if commands.Of(nodeContext.Node.Value) == commands.Run {
						trimStart := len(nodeContext.Node.Value) + 1 // command plus trailing space
						commandText := nodeContext.Node.Original[trimStart:]
						posixCommands, err := shell.NewPosixCommand(commandText)
						if err != nil {
							log.Warnf("Unable to parse RUN command, validation not evaluated for: %s", commandText)
						} else {
							for _, tool := range buildTools {
								re := regexp.MustCompile(tool)
								for _, command := range posixCommands {
									if re.MatchString(command.Name) {
										nodeContext.Context.CausedFailure = true
										hasFailures = true
									}
								}
							}
						}

						validationContexts = append(validationContexts, nodeContext.Context)
					}
				}

				var result model.Valid
				var details string
				if !hasFailures {
					result = model.Success
					details = "No suggestions for multistage builds found."
				} else {
					if hasAnyBuilder {
						result = model.Failure
						details = "Consider moving code compilation tasks into your builder stage."
					} else {
						result = model.Failure
						details = mcr.GetSummary()
					}
				}

				return &validations.ValidationResult{
					Result:   result,
					Details:  details,
					Contexts: validationContexts,
				}
			},
		},
	}

	return &r
}

func init() {
	AddRule(considerMultistageBuild())
}
