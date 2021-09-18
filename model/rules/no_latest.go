package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func isLatest(image string) bool {
	if strings.Contains(image, "@sha256:") {
		return false
	}
	imageParts := strings.SplitAfter(image, ":")
	// FROM scratch isn't considered "latest"
	return (len(imageParts) == 1 && imageParts[0] != "scratch") || imageParts[len(imageParts)-1] == "latest"
}

func processFrom(node *parser.Node, handler func(image string, builderName *string) *validations.ValidationResult) *validations.ValidationResult {
	var image string
	var isBuilder = false
	var builderName *string = nil

	part := node.Next
	for part != nil {
		if part.Value == "as" {
			isBuilder = true
		} else {
			if isBuilder {
				builderName = &part.Value
				break
			}
			image = part.Value
		}
		part = part.Next
	}

	return handler(image, builderName)
}

func validateIfLatest(image string, validationContext validations.ValidationContext, summary string) *validations.ValidationResult {
	if isLatest(image) {
		validationContext.CausedFailure = true
		return &validations.ValidationResult{
			Result:   model.Failure,
			Details:  summary,
			Contexts: []validations.ValidationContext{validationContext},
		}
	}

	return &validations.ValidationResult{
		Result:   model.Success,
		Details:  summary,
		Contexts: []validations.ValidationContext{validationContext},
	}
}

func noLatestBuilder() validations.Rule {
	targetCommands := []commands.DockerCommand{commands.From}
	summary := "Avoid using images tagged as Latest in builder stages"
	rule := validations.SimpleRule{
		Name:     "tagged-latest-builder",
		Summary:  summary,
		Details:  "Using `latest` images in builders is not recommended (builds are not repeatable).",
		Priority: model.LowPriority,
		Commands: targetCommands,
		Handler: func(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
			return processFrom(node, func(image string, builderName *string) *validations.ValidationResult {
				if builderName == nil {
					return validations.NewValidationResultSkipped("No builder reference found in the Dockerfile")
				}
				return validateIfLatest(image, validationContext, summary)
			})
		},

		URL: model.StringPtr("https://docs.docker.com/develop/dev-best-practices/"),
	}
	return &rule
}

func noLatest() validations.Rule {
	targetCommands := []commands.DockerCommand{commands.From}
	summary := "Avoid using images tagged as Latest in production builds"
	rule := validations.SimpleRule{
		Name:     "tagged-latest",
		Summary:  summary,
		Details:  "Docker best practices suggest avoiding `latest` images in production builds",
		Priority: model.HighPriority,
		Commands: targetCommands,
		Handler: func(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
			return processFrom(node, func(image string, builderName *string) *validations.ValidationResult {
				if builderName != nil {
					return validations.NewValidationResultSkipped("This rule does not apply to staged builds")
				}
				return validateIfLatest(image, validationContext, summary)
			})
		},
		URL: model.StringPtr("https://docs.docker.com/develop/dev-best-practices/"),
	}
	return &rule
}

func init() {
	AddRule(noLatest())
	AddRule(noLatestBuilder())
}
