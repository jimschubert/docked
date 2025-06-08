package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func avoidAddExternal() validations.Rule {
	summary := "Avoid using ADD with external files or archives. Use COPY instead."
	r := validations.MultiContextRule{
		Name:    "avoid-add-external",
		Summary: summary,
		Details: "The ADD command supports pulling files over HTTP(s), and auto-extracts some archives. " +
			"Docker's own best practices strongly encourage using COPY of a local file.",
		Priority:         model.CriticalPriority,
		Commands:         []commands.DockerCommand{commands.Add},
		AppliesToBuilder: false,
		URL:              model.StringPtr("https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#add-or-copy"),
		Evaluator: validations.MultiContextPerNodeEvaluator{
			Fn: func(node *parser.Node, validationContext validations.ValidationContext) model.Valid {
				parsed, err := instructions.ParseInstruction(node)
				if err != nil {
					return model.Skipped
				}

				result := model.Success
				if add, ok := parsed.(*instructions.AddCommand); ok && add != nil && len(add.SourcePaths) > 0 {
					input := add.SourcePaths[0]
					if strings.HasPrefix(input, "http:") || strings.HasPrefix(input, "https:") || strings.HasPrefix(input, "file:") {
						result = model.Failure
					} else if !strings.HasSuffix(input, ".tar.xz") {
						result = model.Recommendation
					}
				}

				return result
			},
		},
	}

	return &r
}

func init() {
	AddRule(avoidAddExternal())
}
