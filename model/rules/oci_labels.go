package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

func openContainersAnnotations() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "oci-labels",
		Summary:  "Consider using common annotations defined by Open Containers Initiative",
		Details:  "Open Containers Initiative defines a common set of annotations which expose as labels on containers",
		Priority: model.MediumPriority,
		Commands: []commands.DockerCommand{commands.Label},
		Evaluator: validations.MultiContextFullEvaluator{
			Fn: func(mcr *validations.MultiContextRule) *validations.ValidationResult {
				if mcr == nil || mcr.ContextCache == nil {
					return &validations.ValidationResult{
						Result:  model.Skipped,
						Details: mcr.GetSummary(),
					}
				}
				result := model.Recommendation
				validationContexts := make([]validations.ValidationContext, 0)
				for _, nodeValidationContext := range *mcr.ContextCache {
					if strings.Contains(nodeValidationContext.Node.Original, "org.opencontainers") {
						result = model.Success
					} else {
						nodeValidationContext.Context.HasRecommendations = true
					}
					validationContexts = append(validationContexts, nodeValidationContext.Context)
				}

				return &validations.ValidationResult{
					Result:   result,
					Details:  mcr.GetSummary(),
					Contexts: validationContexts,
				}
			},
		},
	}
	return &r
}

func init() {
	AddRule(openContainersAnnotations())
}
