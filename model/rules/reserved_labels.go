package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"
)

func reservedLabels() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "reserved-labels",
		Summary:  "You can't define labels which are reserved by docker.",
		Details:  "Docker reserves the following namespaces in labels: `com.docker.*`, `io.docker.*`, and `org.dockerproject.*`.",
		Priority: model.CriticalPriority,
		Commands: []commands.DockerCommand{commands.Label},
		Evaluator: validations.MultiContextFullEvaluator{
			Fn: func(mcr *validations.MultiContextRule) *validations.ValidationResult {
				if mcr == nil || mcr.ContextCache == nil {
					return &validations.ValidationResult{
						Result:  model.Skipped,
						Details: mcr.GetSummary(),
					}
				}

				result := model.Success
				validationContexts := make([]validations.ValidationContext, 0)

				for _, nodeValidationContext := range *mcr.ContextCache {
					which, err := instructions.ParseCommand(&nodeValidationContext.Node)
					if err != nil {
						log.Warnf("unabel to parse label: %s", nodeValidationContext.Node.Value)
						continue
					}

					if labelCommand, ok := which.(*instructions.LabelCommand); ok && labelCommand != nil {
						for _, label := range labelCommand.Labels {
							if strings.Contains(label.Key, "com.docker.") ||
								strings.Contains(label.Key, "io.docker.") ||
								strings.Contains(label.Key, "org.dockerproject.") {
								result = model.Failure
								nodeValidationContext.Context.CausedFailure = true
							}
						}
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
	AddRule(reservedLabels())
}
