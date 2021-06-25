package validations

import (
	"fmt"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type EvaluationFunc func(node *parser.Node, validationContext ValidationContext) *ValidationResult
type FinalizeFunc func() *ValidationResult
type ResetFunc func()

type Rule struct {
	Name     string
	Summary  string
	Details  string
	Priority model.Priority
	Commands []commands.DockerCommand
	Evaluate EvaluationFunc
	Category *string
	Finalize *FinalizeFunc
	Reset    *ResetFunc
	URL      *string
}

func (r *Rule) InvokeReset() {
	if r.Reset != nil {
		reset := *r.Reset
		reset()
	}
}

func (r *Rule) InvokeFinalize() *ValidationResult {
	if r.Finalize != nil {
		finalizer := *r.Finalize
		return finalizer()
	}

	return nil
}

func (r *Rule) categoryID() string {
	if r.Category != nil {
		return *r.Category
	}

	if len(r.Commands) == 0 {
		return ""
	}

	switch r.Commands[0] {
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

func (r *Rule) LintID() string {
	return fmt.Sprintf("D%s:%s", r.categoryID(), r.Name)
}