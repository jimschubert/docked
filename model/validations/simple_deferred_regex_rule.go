package validations

import (
	"fmt"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// SimpleDeferredRegexRule is a no-frills regex evaluation which occurs after all relevant nodes of the Dockerfile are parsed and evaluated.
type SimpleDeferredRegexRule struct {
	Name             string                   `json:"name,omitempty"`
	Summary          string                   `json:"summary,omitempty"`
	Details          string                   `json:"details,omitempty"`
	Patterns         []string                 `json:"patterns,omitempty"`
	Priority         model.Priority           `json:"priority,omitempty"`
	Commands         []commands.DockerCommand `json:"commands,omitempty"`
	AppliesToBuilder bool                     `json:"applies_to_builder,omitempty"`
	Category         *string                  `json:"category,omitempty"`
	URL              *string                  `json:"url,omitempty"`
	inBuilderImage   bool
	inFinalImage     bool
	contextCache     *[]NodeValidationContext
}

// GetName gets the name of the rule
func (r *SimpleDeferredRegexRule) GetName() string {
	return r.Name
}

// GetSummary gets the summary of the rule
func (r *SimpleDeferredRegexRule) GetSummary() string {
	return r.Summary
}

// GetDetails gets the details of the rule
func (r *SimpleDeferredRegexRule) GetDetails() string {
	prefix := ""
	if r.Details != "" {
		prefix = fmt.Sprintf("%s\n", r.Details)
	}
	return fmt.Sprintf("%sThis rule matches against the pattern `%s`", prefix, r.Patterns)
}

// GetPriority gets the priority of the rule
func (r *SimpleDeferredRegexRule) GetPriority() model.Priority {
	return r.Priority
}

// GetCommands gets the commands of the rule
func (r *SimpleDeferredRegexRule) GetCommands() []commands.DockerCommand {
	return r.Commands
}

// GetCategory gets the category of the rule
func (r *SimpleDeferredRegexRule) GetCategory() *string {
	return r.Category
}

// GetURL gets the URL of the rule
func (r *SimpleDeferredRegexRule) GetURL() *string {
	return r.URL
}

// GetLintID gets the lint ID of the rule
func (r *SimpleDeferredRegexRule) GetLintID() string {
	return LintID(r)
}

// Evaluate a parsed node and its context
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
		validationContext.IsBuilderContext = true
		if r.AppliesToBuilder {
			*r.contextCache = append(*r.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
		}
	} else {
		*r.contextCache = append(*r.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
	}
	return nil
}

// Reset the rule's internal state
func (r *SimpleDeferredRegexRule) Reset() {
	newCache := make([]NodeValidationContext, 0)
	r.contextCache = &newCache
	r.inBuilderImage = false
	r.inFinalImage = false
}

// Finalize the validation evaluation
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
