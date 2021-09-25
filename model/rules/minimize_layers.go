package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

func minimizeLayers() validations.Rule {
	r := validations.MultiContextRule{
		Name:             "minimize-layers",
		Summary:          "Try to minimize the number of layers which increase image size",
		Details:          "RUN, ADD, and COPY create new layers which may increase the size of the final image. Consider condensing these to fewer than 7 combined layers or use multi-stage builds where possible.",
		Priority:         model.LowPriority,
		Commands:         []commands.DockerCommand{commands.Run, commands.Add, commands.Copy},
		URL:              model.StringPtr("https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#minimize-the-number-of-layers"),
		AppliesToBuilder: false,
		Evaluator: validations.MultiContextFullEvaluator{
			Fn: func(mcr *validations.MultiContextRule) *validations.ValidationResult {
				if mcr == nil || mcr.ContextCache == nil {
					return &validations.ValidationResult{
						Result:  model.Skipped,
						Details: mcr.GetSummary(),
					}
				}

				var result model.Valid
				if len(*mcr.ContextCache) > (len(mcr.Commands) * 2) {
					result = model.Recommendation
				} else {
					result = model.Success
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
	AddRule(minimizeLayers())
}
