package validations

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type MultiContextRule struct {
	Name             string
	Summary          string
	Details          string
	Priority         model.Priority
	Commands         []commands.DockerCommand
	AppliesToBuilder bool
	Category         *string
	URL              *string
	Evaluator        func(node *parser.Node, validationContext ValidationContext) model.Valid
	inBuilderImage   bool
	inFinalImage     bool
	contextCache     *[]NodeValidationContext
}

func (m *MultiContextRule) GetName() string {
	return m.Name
}

func (m *MultiContextRule) GetSummary() string {
	return m.Summary
}

func (m *MultiContextRule) GetDetails() string {
	return m.Details
}

func (m *MultiContextRule) GetPriority() model.Priority {
	return m.Priority
}

func (m *MultiContextRule) GetCommands() []commands.DockerCommand {
	return m.Commands
}

func (m *MultiContextRule) GetCategory() *string {
	return m.Category
}

func (m *MultiContextRule) GetURL() *string {
	return m.URL
}

func (m *MultiContextRule) GetLintID() string {
	return LintID(m)
}

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

func (m *MultiContextRule) Reset() {
	newCache := make([]NodeValidationContext, 0)
	m.contextCache = &newCache
	m.inBuilderImage = false
	m.inFinalImage = false
}

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
