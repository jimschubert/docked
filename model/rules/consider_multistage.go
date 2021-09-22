package rules

import (
	"regexp"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

func considerMultistageBuild() validations.Rule {
	var buildTools = [...]string{`np[mx] install`, `mvn[w]?[ ]`, `bazel build`, `gradle[w]?[ ]`, `\bgo build\b`, `\bgoreleaser\b`}
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
					if nodeContext.Node.Value == string(commands.Run) {
						for _, tool := range buildTools {
							re := regexp.MustCompile(tool)
							if re.MatchString(nodeContext.Node.Original) {
								nodeContext.Context.CausedFailure = true
								hasFailures = true
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
