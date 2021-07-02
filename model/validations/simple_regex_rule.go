package validations

import (
	"fmt"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type SimpleRegexRule struct {
	name      string
	summary   string
	pattern   string
	priority  model.Priority
	command   commands.DockerCommand
	category  *string
	url       *string
	_commands []commands.DockerCommand
}

func (r SimpleRegexRule) Name() string {
	return r.name
}

func (r SimpleRegexRule) Summary() string {
	return r.summary
}

func (r SimpleRegexRule) Details() string {
	return fmt.Sprintf("Found a string matching the pattern %s", r.pattern)
}

func (r SimpleRegexRule) Priority() model.Priority {
	return r.priority
}

func (r SimpleRegexRule) Commands() []commands.DockerCommand {
	if len(r._commands) > 0 {
		return r._commands
	}

	r._commands = append(r._commands, r.command)
	return r._commands
}

func (r SimpleRegexRule) Category() *string {
	return r.category
}

func (r SimpleRegexRule) URL() *string {
	return r.url
}

func (r SimpleRegexRule) LintID() string {
	return LintID(r)
}

func (r SimpleRegexRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	trimStart := len(node.Value) + 1 // command plus trailing space
	matchAgainst := node.Original[trimStart:]
	if model.NewPattern(r.pattern).Matches(matchAgainst) {
		validationContext.CausedFailure = true
		return &ValidationResult{
			Result:   model.Failure,
			Details:  r.Summary(),
			Contexts: []ValidationContext{validationContext},
		}
	}

	return &ValidationResult{
		Result:   model.Success,
		Details:  r.Summary(),
		Contexts: []ValidationContext{validationContext},
	}
}

func NewSimpleRegexRule(
	name string,
	summary string,
	pattern string,
	priority model.Priority,
	command commands.DockerCommand,
	category *string,
	url *string,
) Rule {
	return SimpleRegexRule{
		name:     name,
		summary:  summary,
		pattern:  pattern,
		priority: priority,
		command:  command,
		category: category,
		url:      url,
	}
}
