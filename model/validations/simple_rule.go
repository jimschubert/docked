package validations

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// handlerFunc is a type alias for evaluating a SimpleRule
type handlerFunc func(node *parser.Node, validationContext ValidationContext) *ValidationResult

// SimpleRule is the simplest implementation of a rule
type SimpleRule struct {
	Name     string                   `json:"name,omitempty"`
	Summary  string                   `json:"summary,omitempty"`
	Details  string                   `json:"details,omitempty"`
	Priority model.Priority           `json:"priority,omitempty"`
	Commands []commands.DockerCommand `json:"commands,omitempty"`
	Handler  handlerFunc              `json:"-"`
	Category *string                  `json:"category,omitempty"`
	URL      *string                  `json:"url,omitempty"`
}

// GetName gets the name of the rule
func (r SimpleRule) GetName() string {
	return r.Name
}

// GetSummary gets the summary of the rule
func (r SimpleRule) GetSummary() string {
	return r.Summary
}

// GetDetails gets the details of the rule
func (r SimpleRule) GetDetails() string {
	return r.Details
}

// GetPriority gets the priority of the rule
func (r SimpleRule) GetPriority() model.Priority {
	return r.Priority
}

// GetCommands gets the commands of the rule
func (r SimpleRule) GetCommands() []commands.DockerCommand {
	return r.Commands
}

// GetCategory gets the category of the rule
func (r SimpleRule) GetCategory() *string {
	return r.Category
}

// GetURL gets the URL of the rule
func (r SimpleRule) GetURL() *string {
	return r.URL
}

// GetLintID gets the lint ID of the rule
func (r SimpleRule) GetLintID() string {
	return LintID(r)
}

// Evaluate a parsed node and its context
func (r SimpleRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	return r.Handler(node, validationContext)
}
