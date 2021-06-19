package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/command"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func noLatest() Rule {
	return Rule{
		Name: "tagged-latest",
		Command: command.From,
		Evaluate: func(node *parser.Node) *validations.ValidationResult {
			var ruleName string
			var priority model.Priority
			var details string
			var image string
			var isBuilder = false

			part := node.Next
			for part != nil {
				if part.Value == "as" {
					isBuilder = true
					break
				} else {
					image = part.Value
				}
				part = part.Next
			}

			if isBuilder {
				ruleName = "latest-builder"
				priority = model.LowPriority
				details = "Using 'latest' images in builders is not recommended."
			} else {
				ruleName = "latest-final"
				priority = model.HighPriority
				details = "Docker best practices suggest avoiding 'latest' images in production builds"
			}

			imageParts := strings.Split(image, ":")
			if len(imageParts) == 1 || imageParts[1] == "latest" {
				return validations.NewFailureResult(ruleName, priority, details)
			}

			return validations.NewSuccessResult(ruleName, priority, details)
		},
	}
}

func init() {
	AddRule(noLatest())
}
