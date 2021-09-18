package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type singleCmd struct {
	contextCache *[]validations.NodeValidationContext
}

func (s *singleCmd) GetName() string {
	return "single-cmd"
}

func (s *singleCmd) GetSummary() string {
	return "Only a single CMD instruction is supported"
}

func (s *singleCmd) GetDetails() string {
	return "More than one CMD may indicate a programming error. Docker will run the last CMD instruction only, but this could be a security concern."
}

func (s *singleCmd) GetPriority() model.Priority {
	return model.CriticalPriority
}

func (s *singleCmd) GetCommands() []commands.DockerCommand {
	return []commands.DockerCommand{commands.Cmd}
}

func (s *singleCmd) GetCategory() *string {
	return nil
}

func (s *singleCmd) GetURL() *string {
	return model.StringPtr("https://docs.docker.com/engine/reference/builder/#cmd")
}

func (s *singleCmd) GetLintID() string {
	return validations.LintID(s)
}

func (s *singleCmd) Evaluate(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
	*s.contextCache = append(*s.contextCache, validations.NodeValidationContext{Node: *node, Context: validationContext})
	return nil
}

func (s *singleCmd) Reset() {
	newCache := make([]validations.NodeValidationContext, 0)
	s.contextCache = &newCache
}

func (s *singleCmd) Finalize() *validations.ValidationResult {
	validationContexts := make([]validations.ValidationContext, 0)
	for _, nodeContext := range *s.contextCache {
		validationContexts = append(validationContexts, nodeContext.Context)
	}

	var result model.Valid
	if len(*s.contextCache) > 1 {
		result = model.Failure
	} else {
		result = model.Success
	}

	return &validations.ValidationResult{
		Result:   result,
		Details:  s.GetSummary(),
		Contexts: validationContexts,
	}
}

func init() {
	singleCmdRule := singleCmd{}
	AddRule(&singleCmdRule)
}
