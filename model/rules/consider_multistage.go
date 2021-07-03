package rules

import (
	"regexp"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

var buildTools = [...]string{`np[mx] install`, `mvn[w]?[ ]`, `bazel build`, `gradle[w]?[ ]`}

type considerMultistageBuild struct {
	contextCache   *[]validations.NodeValidationContext
	hasAnyBuilder  bool
	inFinalContext bool
}

func (c *considerMultistageBuild) Name() string {
	return "consider-multistage"
}

func (c *considerMultistageBuild) Summary() string {
	return "Consider using multi-stage builds for complex operations like building code."
}

func (c *considerMultistageBuild) Details() string {
	return "A multi-stage build can reduce the final image size by building necessary components or" +
		" downloading large archives in a separate build context. This can help keep your final image lean."
}

func (c *considerMultistageBuild) Priority() model.Priority {
	return model.LowPriority
}

func (c *considerMultistageBuild) Commands() []commands.DockerCommand {
	return []commands.DockerCommand{commands.Run, commands.From}
}

func (c *considerMultistageBuild) Category() *string {
	return nil
}

func (c *considerMultistageBuild) URL() *string {
	return model.StringPtr("https://docs.docker.com/develop/develop-images/multistage-build/")
}

func (c *considerMultistageBuild) LintID() string {
	return validations.LintID(c)
}

func (c *considerMultistageBuild) Evaluate(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
	if !c.hasAnyBuilder {
		c.hasAnyBuilder = model.IsBuilderFrom(node)
	}
	if !c.inFinalContext {
		c.inFinalContext = model.IsFinalFrom(node)
	}

	if c.inFinalContext {
		*c.contextCache = append(*c.contextCache, validations.NodeValidationContext{Node: *node, Context: validationContext})
	}
	return nil
}

func (c *considerMultistageBuild) Reset() {
	newCache := make([]validations.NodeValidationContext, 0)
	c.contextCache = &newCache
	c.hasAnyBuilder = false
	c.inFinalContext = false
}

func (c *considerMultistageBuild) Finalize() *validations.ValidationResult {
	var hasFailures bool
	validationContexts := make([]validations.ValidationContext, 0)
	for _, nodeContext := range *c.contextCache {
		validationContexts = append(validationContexts, nodeContext.Context)
		if nodeContext.Node.Value == string(commands.Run) {
			for _, tool := range buildTools {
				re := regexp.MustCompile(tool)
				if re.MatchString(nodeContext.Node.Original) {
					nodeContext.Context.CausedFailure = true
					hasFailures = true
				}
			}
		}
	}

	var result model.Valid
	var details string
	if !hasFailures {
		result = model.Success
		details = "No suggestions for multistage builds found."
	} else {
		if c.hasAnyBuilder {
			result = model.Failure
			details = "Consider moving code compilation tasks into your builder stage."
		} else {
			result = model.Failure
			details = c.Summary()
		}
	}

	return &validations.ValidationResult{
		Result:   result,
		Details:  details,
		Contexts: validationContexts,
	}
}

func init() {
	consideration := considerMultistageBuild{}
	AddRule(&consideration)
}
