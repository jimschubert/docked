package validations

import (
	"fmt"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type SimpleDeferredRegexRule struct {
	Name             string
	Summary          string
	Details          string
	Patterns         []string
	Priority         model.Priority
	Commands         []commands.DockerCommand
	AppliesToBuilder bool
	Category         *string
	URL              *string
	inBuilderImage   bool
	inFinalImage     bool
	contextCache     *[]NodeValidationContext
}

func (r *SimpleDeferredRegexRule) GetName() string {
	return r.Name
}

func (r *SimpleDeferredRegexRule) GetSummary() string {
	return r.Summary
}

func (r *SimpleDeferredRegexRule) GetDetails() string {
	prefix := ""
	if r.Details != "" {
		prefix = fmt.Sprintf("%s\n", r.Details)
	}
	return fmt.Sprintf("%sThis rule matches against the pattern `%s`", prefix, r.Patterns)
}

func (r *SimpleDeferredRegexRule) GetPriority() model.Priority {
	return r.Priority
}

func (r *SimpleDeferredRegexRule) GetCommands() []commands.DockerCommand {
	return r.Commands
}

func (r *SimpleDeferredRegexRule) GetCategory() *string {
	return r.Category
}

func (r *SimpleDeferredRegexRule) GetURL() *string {
	return r.URL
}

func (r *SimpleDeferredRegexRule) GetLintID() string {
	return LintID(r)
}

func (r *SimpleDeferredRegexRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	if !r.inBuilderImage {
		r.inBuilderImage = model.IsBuilderFrom(node)
	}
	if !r.inFinalImage {
		r.inFinalImage = model.IsFinalFrom(node)
	}
	if r.inFinalImage && r.inBuilderImage {
		r.inBuilderImage = false
	}

	if r.inBuilderImage {
		if r.AppliesToBuilder {
			*r.contextCache = append(*r.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
		}
	} else {
		*r.contextCache = append(*r.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
	}
	return nil
}

func (r *SimpleDeferredRegexRule) Reset() {
	newCache := make([]NodeValidationContext, 0)
	r.contextCache = &newCache
	r.inBuilderImage = false
	r.inFinalImage = false
}

func (r *SimpleDeferredRegexRule) Finalize() *ValidationResult {
	validationContexts := make([]ValidationContext, 0)
	hasFailures := false
	for _, nodeContext := range *r.contextCache {
		trimStart := len(nodeContext.Node.Value) + 1 // command plus trailing space
		matchAgainst := nodeContext.Node.Original[trimStart:]
		for _, pattern := range r.Patterns {
			if model.NewPattern(pattern).Matches(matchAgainst) {
				nodeContext.Context.CausedFailure = true
				hasFailures = true
			}
		}
		validationContexts = append(validationContexts, nodeContext.Context)
	}

	if hasFailures {
		return &ValidationResult{
			Result:   model.Failure,
			Details:  r.GetSummary(),
			Contexts: validationContexts,
		}
	}

	return &ValidationResult{
		Result:   model.Success,
		Details:  r.GetSummary(),
		Contexts: validationContexts,
	}
}
