package validations

import (
	"fmt"
	"regexp"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type SimpleDeferredRegexRule struct {
	name             string
	summary          string
	patterns         []string
	priority         model.Priority
	commands         []commands.DockerCommand
	appliesToBuilder bool
	category         *string
	url              *string
	inBuilderImage   bool
	inFinalImage     bool
	contextCache     *[]NodeValidationContext
}

func (r SimpleDeferredRegexRule) Name() string {
	return r.name
}

func (r SimpleDeferredRegexRule) Summary() string {
	return r.summary
}

func (r SimpleDeferredRegexRule) Details() string {
	return fmt.Sprintf("Found a string matching %s", r.patterns)
}

func (r SimpleDeferredRegexRule) Priority() model.Priority {
	return r.priority
}

func (r SimpleDeferredRegexRule) Commands() []commands.DockerCommand {
	return r.commands
}

func (r SimpleDeferredRegexRule) Category() *string {
	return r.category
}

func (r SimpleDeferredRegexRule) URL() *string {
	return r.url
}

func (r SimpleDeferredRegexRule) LintID() string {
	return LintID(r)
}

func (r SimpleDeferredRegexRule) Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult {
	if !r.inBuilderImage {
		r.inBuilderImage = model.IsBuilderFrom(node)
	}
	if !r.inFinalImage {
		r.inFinalImage = model.IsFinalFrom(node)
	}
	if r.inFinalImage && r.inBuilderImage {
		r.inBuilderImage = false
	}
	if r.inBuilderImage && r.appliesToBuilder {
		*r.contextCache = append(*r.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
	} else {
		*r.contextCache = append(*r.contextCache, NodeValidationContext{Node: *node, Context: validationContext})
	}
	return nil
}

func (r SimpleDeferredRegexRule) Reset() {
	newCache := make([]NodeValidationContext, 0)
	r.contextCache = &newCache
	r.inBuilderImage = false
	r.inFinalImage = false
}

func (r SimpleDeferredRegexRule) Finalize() *ValidationResult {
	validationContexts := make([]ValidationContext, 0)
	for _, nodeContext := range *r.contextCache {
		trimStart := len(nodeContext.Node.Value) + 1 // command plus trailing space
		matchAgainst := nodeContext.Node.Original[trimStart:]
		for _, pattern := range r.patterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(matchAgainst) {
				return &ValidationResult{
					Result:   model.Failure,
					Details:  r.Summary(),
					Contexts: validationContexts,
				}
			}
		}
	}

	return &ValidationResult{
		Result:   model.Success,
		Details:  r.Summary(),
		Contexts: validationContexts,
	}
}

func NewSimpleDeferredRegexRule(
	name string,
	summary string,
	patterns []string,
	priority model.Priority,
	commands []commands.DockerCommand,
	appliesToBuilder bool,
	category *string,
	url *string,
) Rule {
	return SimpleDeferredRegexRule{
		name:             name,
		summary:          summary,
		patterns:         patterns,
		priority:         priority,
		commands:         commands,
		appliesToBuilder: appliesToBuilder,
		category:         category,
		url:              url,
	}
}
