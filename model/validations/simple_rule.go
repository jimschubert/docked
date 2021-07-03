package validations

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type SimpleRule struct {
	Name     string
	Summary  string
	Details  string
	Priority model.Priority
	Commands []commands.DockerCommand
	Handler  func(node *parser.Node, validationContext ValidationContext) *ValidationResult
	Category *string
	URL      *string
}

func (r SimpleRule) GetName() string {
	return r.Name
}

func (r SimpleRule) GetSummary() string {
	return r.Summary
}

func (r SimpleRule) GetDetails() string {
	return r.Details
}

func (r SimpleRule) GetPriority() model.Priority {
	return r.Priority
}

func (r SimpleRule) GetCommands() []commands.DockerCommand {
	return r.Commands
}

func (r SimpleRule) GetCategory() *string {
	return r.Category
}

func (r SimpleRule) GetURL() *string {
	return r.URL
}

func (r SimpleRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	return r.Handler(node, validationContext)
}

func (r SimpleRule) GetLintID() string {
	return LintID(r)
}
