package validations

import (
	"fmt"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type SimpleRegexRule struct {
	Name      string                 `json:"name,omitempty"`
	Summary   string                 `json:"summary,omitempty"`
	Details   string                 `json:"details,omitempty"`
	Pattern   string                 `json:"pattern,omitempty"`
	Priority  model.Priority         `json:"priority,omitempty"`
	Command   commands.DockerCommand `json:"command,omitempty"`
	Category  *string                `json:"category,omitempty"`
	URL       *string                `json:"url,omitempty"`
	_commands []commands.DockerCommand
}

func (r SimpleRegexRule) GetName() string {
	return r.Name
}

func (r SimpleRegexRule) GetSummary() string {
	return r.Summary
}

func (r SimpleRegexRule) GetDetails() string {
	prefix := ""
	if r.Details != "" {
		prefix = fmt.Sprintf("%s\n", r.Details)
	}
	return fmt.Sprintf("%sThis rule matches against the pattern `%s`", prefix, r.Pattern)
}

func (r SimpleRegexRule) GetPriority() model.Priority {
	return r.Priority
}

func (r SimpleRegexRule) GetCommands() []commands.DockerCommand {
	if len(r._commands) > 0 {
		return r._commands
	}

	r._commands = append(r._commands, r.Command)
	return r._commands
}

func (r SimpleRegexRule) GetCategory() *string {
	return r.Category
}

func (r SimpleRegexRule) GetURL() *string {
	return r.URL
}

func (r SimpleRegexRule) GetLintID() string {
	return LintID(r)
}

func (r SimpleRegexRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	trimStart := len(node.Value) + 1 // command plus trailing space
	matchAgainst := node.Original[trimStart:]
	if model.NewPattern(r.Pattern).Matches(matchAgainst) {
		validationContext.CausedFailure = true
		return &ValidationResult{
			Result:   model.Failure,
			Details:  r.GetSummary(),
			Contexts: []ValidationContext{validationContext},
		}
	}

	return &ValidationResult{
		Result:   model.Success,
		Details:  r.GetSummary(),
		Contexts: []ValidationContext{validationContext},
	}
}
