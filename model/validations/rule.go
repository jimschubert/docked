package validations

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// ResettingRule defines the behaviors for a rule which can have its state externally reset.
// It is up to values implementing this interface to appropriately lock resources as needed.
type ResettingRule interface {
	Rule
	// Reset the rule's internal state
	Reset()
}

// FinalizingRule defines the behaviors for a rule which performs optional post-processing or finalization before returning a ValidationResult
type FinalizingRule interface {
	ResettingRule
	// Finalize the validation evaluation
	Finalize() *ValidationResult
}

// Rule defines the immutable interface and behaviors for those types implementing evaluations and their details
type Rule interface {
	// GetName gets the name of the rule
	GetName() string
	// GetSummary gets the summary of the rule
	GetSummary() string
	// GetDetails gets the details of the rule
	GetDetails() string
	// GetPriority gets the priority of the rule
	GetPriority() model.Priority
	// GetCommands gets the commands of the rule
	GetCommands() []commands.DockerCommand
	// GetCategory gets the category of the rule
	GetCategory() *string
	// GetURL gets the URL of the rule
	GetURL() *string
	// GetLintID gets the lint ID of the rule
	GetLintID() string
	// Evaluate a parsed node and its context
	Evaluate(node *parser.Node, validationContext ValidationContext) *ValidationResult
}

// LintID builds out a standard format for identifying a linting rule
// A rule's name is converted to alphanumeric kebab-case here.
func LintID(rule Rule) string {
	buf := bytes.Buffer{}
	name := strings.ToLower(rule.GetName())
	for _, r := range name {
		if unicode.IsSpace(r) || r == '-' {
			buf.WriteRune('-')
		} else if unicode.IsNumber(r) || unicode.IsLetter(r) {
			buf.WriteRune(r)
		}
	}
	return fmt.Sprintf("D%s:%s", CategoryID(rule), buf.String())
}

// CategoryID determines an identifier in relation to the rule
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
	return commandCategoryCharacter(cmd)
}

func commandCategoryCharacter(cmd commands.DockerCommand) string {
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
