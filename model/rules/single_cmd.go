package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

func singleCmd() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "single-cmd",
		Summary:  "Only a single CMD instruction is supported",
		Details:  "More than one CMD may indicate a programming error. Docker will run the last CMD instruction only, but this could be a security concern.",
		Priority: model.CriticalPriority,
		Commands: []commands.DockerCommand{commands.Cmd},
		URL:      model.StringPtr("https://docs.docker.com/engine/reference/builder/#cmd"),
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
				for _, nodeContext := range *mcr.ContextCache {
					result = model.Failure
					validationContexts = append(validationContexts, nodeContext.Context)
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
	AddRule(singleCmd())
}
