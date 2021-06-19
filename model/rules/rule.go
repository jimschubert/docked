package rules

import (
	"fmt"

	"github.com/jimschubert/docked/model/docker/command"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type EvaluationFunc func(node *parser.Node) *validations.ValidationResult
type Rule struct {
	Name     string
	Command  command.DockerCommand
	Evaluate EvaluationFunc
}

func (r *Rule) categoryID() string {
	switch r.Command {
	case command.Add:
		return "0"
	case command.Arg:
		return "1"
	case command.Cmd:
		return "2"
	case command.Copy:
		return "3"
	case command.Entrypoint:
		return "4"
	case command.Env:
		return "5"
	case command.Expose:
		return "6"
	case command.From:
		return "7"
	case command.Healthcheck:
		return "8"
	case command.Label:
		return "9"
	case command.Maintainer:
		return "A"
	case command.Onbuild:
		return "B"
	case command.Run:
		return "C"
	case command.Shell:
		return "D"
	case command.StopSignal:
		return "E"
	case command.User:
		return "F"
	case command.Volume:
		return "G"
	case command.Workdir:
		return "H"
	default:
		return ""
	}
}

func (r *Rule) LintID() string {
	return fmt.Sprintf("D%s:%s", r.categoryID(), r.Name)
}
