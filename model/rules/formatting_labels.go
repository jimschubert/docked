package rules

import (
	"strings"
	"unicode"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

func formattingLabels() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "formatting-labels",
		Summary:  "Label keys should be formatted correctly.",
		Details:  "Label keys should begin and end with a lower-case letter and should only contain lower-case alphanumeric characters, the period character (.), and the hyphen character (-). Consecutive periods or hyphens are not allowed.",
		Priority: model.HighPriority,
		Commands: []commands.DockerCommand{commands.Label},
		Evaluator: validations.MultiContextPerNodeEvaluator{
			Fn: func(node *parser.Node, validationContext validations.ValidationContext) model.Valid {
				which, err := instructions.ParseCommand(node)
				if err != nil {
					log.Warnf("unabel to parse label: %s", node.Value)
					return model.Skipped
				}

				if labelCommand, ok := which.(*instructions.LabelCommand); ok && labelCommand != nil {
					illegalRune := func(c rune) bool {
						if unicode.IsLetter(c) {
							return !unicode.IsLower(c)
						}
						if unicode.IsNumber(c) || c == '.' || c == '-' {
							return false
						}
						return true
					}
					for _, label := range labelCommand.Labels {
						illegalRunes := make([]rune, 0)
						for _, c := range label.Key {
							if illegalRune(c) {
								illegalRunes = append(illegalRunes, c)
							}
						}
						if len(illegalRunes) > 0 || strings.Contains(label.Key, "..") || strings.Contains(label.Key, "--") {
							return model.Failure
						}
					}
				}

				return model.Success
			},
		},
	}
	return &r
}

func init() {
	AddRule(formattingLabels())
}
