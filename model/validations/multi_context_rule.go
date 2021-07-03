package validations

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type multiContextRule struct {
	name             string
	summary          string
	details          string
	priority         model.Priority
	commands         []commands.DockerCommand
	appliesToBuilder bool
	category         *string
	url              *string
	inBuilderImage   bool
	inFinalImage     bool
	contextCache     *[]NodeValidationContext
	evaluator        func(node *parser.Node, validationContext ValidationContext) model.Valid
}

func (m *multiContextRule) Name() string {
	return m.name
}

func (m *multiContextRule) Summary() string {
	return m.summary
}

func (m *multiContextRule) Details() string {
	return m.details
}

func (m *multiContextRule) Priority() model.Priority {
	return m.priority
}

func (m *multiContextRule) Commands() []commands.DockerCommand {
	return m.commands
}

func (m *multiContextRule) Category() *string {
	return m.category
}

func (m *multiContextRule) URL() *string {
	return m.url
}

func (m *multiContextRule) LintID() string {
	return LintID(m)
}

func (m *multiContextRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	if !m.inBuilderImage {
		m.inBuilderImage = model.IsBuilderFrom(node)
	}
	if !m.inFinalImage {
		m.inFinalImage = model.IsFinalFrom(node)
	}

	if m.inBuilderImage {
		if m.appliesToBuilder {
			*m.contextCache = append(*m.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
		}
	} else {
		*m.contextCache = append(*m.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
	}
	return nil
}

func (m *multiContextRule) Reset() {
	newCache := make([]NodeValidationContext, 0)
	m.contextCache = &newCache
	m.inBuilderImage = false
	m.inFinalImage = false
}

func (m *multiContextRule) Finalize() *ValidationResult {
	result := model.Success
	validationContexts := make([]ValidationContext, 0)
	for _, nodeContext := range *m.contextCache {
		state := m.evaluator(&nodeContext.Node, nodeContext.Context)
		if state == model.Failure {
			nodeContext.Context.CausedFailure = true
			result = model.Failure
		}
		validationContexts = append(validationContexts, nodeContext.Context)
	}
	return &ValidationResult{
		Result:   result,
		Details:  m.Summary(),
		Contexts: validationContexts,
	}
}

func NewMultiContextRule(
	name string,
	summary string,
	details string,
	priority model.Priority,
	commands []commands.DockerCommand,
	appliesToBuilder bool,
	category *string,
	url *string,
	evaluator func(node *parser.Node, validationContext ValidationContext) model.Valid,
) Rule {
	r := multiContextRule{
		name:             name,
		summary:          summary,
		details:          details,
		priority:         priority,
		commands:         commands,
		appliesToBuilder: appliesToBuilder,
		category:         category,
		url:              url,
		evaluator:        evaluator,
	}
	return &r
}
