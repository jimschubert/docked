package validations

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// MultiContextRule is a rule which spans the context of one or more lines, suggesting a need for deferred evaluation of that rule.
type MultiContextRule struct {
	Name             string                   `json:"name,omitempty"`
	Summary          string                   `json:"summary,omitempty"`
	Details          string                   `json:"details,omitempty"`
	Priority         model.Priority           `json:"priority,omitempty"`
	Commands         []commands.DockerCommand `json:"commands,omitempty"`
	AppliesToBuilder bool                     `json:"applies_to_builder,omitempty"`
	Category         *string                  `json:"category,omitempty"`
	URL              *string                  `json:"url,omitempty"`
	Evaluator        MultiContextEvaluator    `json:"-"`
	ContextCache     *[]NodeValidationContext `json:"-"`
	inBuilderImage   bool
	inFinalImage     bool
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
	if validationContext.IsBuilderContext && m.AppliesToBuilder {
		*m.ContextCache = append(*m.ContextCache, NodeValidationContext{Node: *node, Context: validationContext})
	}
	if !validationContext.IsBuilderContext {
		*m.ContextCache = append(*m.ContextCache, NodeValidationContext{Node: *node, Context: validationContext})
	}
	return nil
}

// Reset the rule's internal state
func (m *MultiContextRule) Reset() {
	newCache := make([]NodeValidationContext, 0)
	m.ContextCache = &newCache
	m.inBuilderImage = false
	m.inFinalImage = false
}

// Finalize the validation evaluation
func (m *MultiContextRule) Finalize() *ValidationResult {
	return m.Evaluator.Evaluate(m)
}

// MultiContextEvaluator defines the Evaluate interface used by MultiContextRule in the Finalize step
type MultiContextEvaluator interface {
	// Evaluate a given MultiContextRule to determine the final ValidationResult
	Evaluate(mcr *MultiContextRule) *ValidationResult
}

// MultiContextFullEvaluator evaluates the MultiContextRule contextually,
// where the handler function has access to the entire NodeValidationContext cache
type MultiContextFullEvaluator struct {
	// Fn evaluates a given MultiContextRule to determine the final ValidationResult
	Fn func(mcr *MultiContextRule) *ValidationResult
}

// Evaluate a given MultiContextRule to determine the final ValidationResult
func (m MultiContextFullEvaluator) Evaluate(mcr *MultiContextRule) *ValidationResult {
	return m.Fn(mcr)
}

// MultiContextPerNodeEvaluator evaluates each node entry against the Handler Fn
type MultiContextPerNodeEvaluator struct {
	// Fn evaluates a parser.Node and its associated ValidationContext to determine if the context is valid
	Fn func(node *parser.Node, validationContext ValidationContext) model.Valid
}

// Evaluate a given MultiContextRule to determine the final ValidationResult
func (m MultiContextPerNodeEvaluator) Evaluate(mcr *MultiContextRule) *ValidationResult {
	result := model.Success
	validationContexts := make([]ValidationContext, 0)
	for _, nodeContext := range *mcr.ContextCache {
		state := m.Fn(&nodeContext.Node, nodeContext.Context)
		if state == model.Failure {
			nodeContext.Context.CausedFailure = true
			result = model.Failure
		}
		if state == model.Recommendation {
			nodeContext.Context.HasRecommendations = true
			result = model.Recommendation
		}
		validationContexts = append(validationContexts, nodeContext.Context)
	}
	return &ValidationResult{
		Result:   result,
		Details:  mcr.GetSummary(),
		Contexts: validationContexts,
	}
}
