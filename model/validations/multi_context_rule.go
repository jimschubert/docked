package validations

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// MultiContextRule is a rule which spans the context of one or more lines, suggesting a need for deferred evaluation of that rule.
type MultiContextRule struct {
	Name             string                                                                   `json:"name,omitempty"`
	Summary          string                                                                   `json:"summary,omitempty"`
	Details          string                                                                   `json:"details,omitempty"`
	Priority         model.Priority                                                           `json:"priority,omitempty"`
	Commands         []commands.DockerCommand                                                 `json:"commands,omitempty"`
	AppliesToBuilder bool                                                                     `json:"applies_to_builder,omitempty"`
	Category         *string                                                                  `json:"category,omitempty"`
	URL              *string                                                                  `json:"url,omitempty"`
	Evaluator        func(node *parser.Node, validationContext ValidationContext) model.Valid `json:"-"`
	inBuilderImage   bool
	inFinalImage     bool
	contextCache     *[]NodeValidationContext
}

// GetName gets the name of the rule
func (m *MultiContextRule) GetName() string {
	return m.Name
}

// GetSummary gets the summary of the rule
func (m *MultiContextRule) GetSummary() string {
	return m.Summary
}

// GetDetails gets the details of the rule
func (m *MultiContextRule) GetDetails() string {
	return m.Details
}

// GetPriority gets the priority of the rule
func (m *MultiContextRule) GetPriority() model.Priority {
	return m.Priority
}

// GetCommands gets the commands of the rule
func (m *MultiContextRule) GetCommands() []commands.DockerCommand {
	return m.Commands
}

// GetCategory gets the category of the rule
func (m *MultiContextRule) GetCategory() *string {
	return m.Category
}

// GetURL gets the URL of the rule
func (m *MultiContextRule) GetURL() *string {
	return m.URL
}

// GetLintID gets the lint ID of the rule
func (m *MultiContextRule) GetLintID() string {
	return LintID(m)
}

// Evaluate a parsed node and its context
func (m *MultiContextRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	if !m.inBuilderImage {
		m.inBuilderImage = model.IsBuilderFrom(node)
	}
	if !m.inFinalImage {
		m.inFinalImage = model.IsFinalFrom(node)
	}

	if m.inBuilderImage {
		if m.AppliesToBuilder {
			*m.contextCache = append(*m.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
		}
	} else {
		*m.contextCache = append(*m.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
	}
	return nil
}

// Reset the rule's internal state
func (m *MultiContextRule) Reset() {
	newCache := make([]NodeValidationContext, 0)
	m.contextCache = &newCache
	m.inBuilderImage = false
	m.inFinalImage = false
}

// Finalize the validation evaluation
func (m *MultiContextRule) Finalize() *ValidationResult {
	result := model.Success
	validationContexts := make([]ValidationContext, 0)
	for _, nodeContext := range *m.contextCache {
		state := m.Evaluator(&nodeContext.Node, nodeContext.Context)
		if state == model.Failure {
			nodeContext.Context.CausedFailure = true
			result = model.Failure
		}
		validationContexts = append(validationContexts, nodeContext.Context)
	}
	return &ValidationResult{
		Result:   result,
		Details:  m.GetSummary(),
		Contexts: validationContexts,
	}
}
