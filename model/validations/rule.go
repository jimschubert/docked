package validations

import (
	"fmt"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type ResettingRule interface {
	Rule
	Reset()
}

type FinalizingRule interface {
	ResettingRule
	Finalize() *ValidationResult
}

type Rule interface {
	GetName() string
	GetSummary() string
	GetDetails() string
	GetPriority() model.Priority
	GetCommands() []commands.DockerCommand
	GetCategory() *string
	GetURL() *string
	GetLintID() string
	Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult
}

func LintID(rule Rule) string {
	return fmt.Sprintf("D%s:%s", CategoryID(rule), rule.GetName())
}

func CategoryID(rule Rule) string {
	category := rule.GetCategory()
	if category != nil {
		return *category
	}
	ruleCommands := rule.GetCommands()
	if len(ruleCommands) == 0 {
		return ""
	}
	cmd := ruleCommands[0]
	return CommandCategoryCharacter(cmd)
}

func CommandCategoryCharacter(cmd commands.DockerCommand) string {
	switch cmd {
	case commands.Add:
		return "0"
	case commands.Arg:
		return "1"
	case commands.Cmd:
		return "2"
	case commands.Copy:
		return "3"
	case commands.Entrypoint:
		return "4"
	case commands.Env:
		return "5"
	case commands.Expose:
		return "6"
	case commands.From:
		return "7"
	case commands.Healthcheck:
		return "8"
	case commands.Label:
		return "9"
	case commands.Maintainer:
		return "A"
	case commands.Onbuild:
		return "B"
	case commands.Run:
		return "C"
	case commands.Shell:
		return "D"
	case commands.StopSignal:
		return "E"
	case commands.User:
		return "F"
	case commands.Volume:
		return "G"
	case commands.Workdir:
		return "H"
	default:
		return ""
	}
}
