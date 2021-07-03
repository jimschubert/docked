package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)


type curlWithoutFail struct{
	contextCache   *[]validations.NodeValidationContext
}

func (c *curlWithoutFail) Name() string {
	return "curl-without-fail"
}

func (c *curlWithoutFail) Summary() string {
	return "Avoid using curl without the silent failing option -f/--fail"
}

func (c *curlWithoutFail) Details() string {
	// accidentally downloading a HTML 404 document is the worst.
	return "Invoking curl without -f/--fail may result in incorrect, missing or stale data, which is a security concern. " +
		"Ignore this rule only if you're handling server errors or verifying file contents separately."
}

func (c *curlWithoutFail) Priority() model.Priority {
	return model.CriticalPriority
}

func (c *curlWithoutFail) Commands() []commands.DockerCommand {
	return []commands.DockerCommand{commands.Run}
}

func (c *curlWithoutFail) Category() *string {
	return nil
}

func (c *curlWithoutFail) URL() *string {
	return model.StringPtr("https://curl.se/docs/faq.html#Why_do_I_get_downloaded_data_eve")
}

func (c *curlWithoutFail) LintID() string {
	return validations.LintID(c)
}

func (c *curlWithoutFail) Evaluate(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
	*c.contextCache = append(*c.contextCache, validations.NodeValidationContext{Node: *node, Context: validationContext})
	return nil
}

func (c *curlWithoutFail) Reset() {
	newCache := make([]validations.NodeValidationContext, 0)
	c.contextCache = &newCache
}

func (c *curlWithoutFail) Finalize() *validations.ValidationResult {
	validationContexts := make([]validations.ValidationContext, 0)
	result := model.Success
	for _, nodeContext := range *c.contextCache {
		trimStart := len(nodeContext.Node.Value) + 1 // command plus trailing space
		matchAgainst := nodeContext.Node.Original[trimStart:]
		if model.NewPattern(`\bcurl\b`).Matches(matchAgainst) {
			if !model.NewPattern(`\b-f\b|\b--fail\b`).Matches(matchAgainst) {
				result = model.Failure
				nodeContext.Context.CausedFailure = true
			}
		}
		validationContexts = append(validationContexts, nodeContext.Context)
	}

	return &validations.ValidationResult{
		Result:   result,
		Details:  c.Summary(),
		Contexts: validationContexts,
	}
}

func init() {
	r := curlWithoutFail{}
	AddRule(&r)
}
