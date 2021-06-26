package validations

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type SimpleRule struct {
	name     string
	summary  string
	details  string
	priority model.Priority
	commands []commands.DockerCommand
	handler  func(node *parser.Node, validationContext ValidationContext) *ValidationResult
	category *string
	url      *string
}

func (r SimpleRule) Name() string {
	return r.name
}

func (r SimpleRule) Summary() string {
	return r.summary
}

func (r SimpleRule) Details() string {
	return r.details
}

func (r SimpleRule) Priority() model.Priority {
	return r.priority
}

func (r SimpleRule) Commands() []commands.DockerCommand {
	return r.commands
}

func (r SimpleRule) Category() *string {
	return r.category
}

func (r SimpleRule) URL() *string {
	return r.url
}

func (r SimpleRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	return r.handler(node, validationContext)
}

func (r SimpleRule) LintID() string {
	return LintID(r)
}

func NewSimpleRule(
	name string,
	summary string,
	details string,
	priority model.Priority,
	commands []commands.DockerCommand,
	handler func(node *parser.Node, validationContext ValidationContext) *ValidationResult,
	category *string,
	url *string,
) Rule {
	return SimpleRule{
		name:     name,
		summary:  summary,
		details:  details,
		priority: priority,
		commands: commands,
		handler:  handler,
		category: category,
		url:      url,
	}
}
