package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func isLatest(image string) bool {
	imageParts := strings.Split(image, ":")
	// FROM scratch isn't considered "latest"
	return (len(imageParts) == 1 && imageParts[0] != "scratch") || imageParts[1] == "latest"
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

func noLatestBuilder() validations.Rule {
	targetCommands := []commands.DockerCommand{commands.From}
	summary := "Avoid using images tagged as Latest in builder stages"
	return validations.NewSimpleRule(
		"tagged-latest-builder",
		summary,
		"Using 'latest' images in builders is not recommended.",
		model.LowPriority,
		targetCommands,
		func(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
			return processFrom(node, func(image string, builderName *string) *validations.ValidationResult {
				if builderName == nil {
					return validations.NewSkippedResult("No builder reference")
				}
				if isLatest(image) {
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
			})
		},
		nil,

		model.StringPtr("https://docs.docker.com/develop/dev-best-practices/"),
	)
}

func noLatest() validations.Rule {
	targetCommands := []commands.DockerCommand{commands.From}
	summary := "Avoid using images tagged as Latest in production builds"
	return validations.NewSimpleRule(
		"tagged-latest",
		summary,
		"Docker best practices suggest avoiding 'latest' images in production builds",
		model.HighPriority,
		targetCommands,
		func(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
			return processFrom(node, func(image string, builderName *string) *validations.ValidationResult {
				if builderName != nil {
					return validations.NewSkippedResult("This rule does not apply to staged builds")
				}
				if isLatest(image) {
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
			})
		},
		nil,
		model.StringPtr("https://docs.docker.com/develop/dev-best-practices/"),
	)
}

func init() {
	AddRule(noLatest())
	AddRule(noLatestBuilder())
}
