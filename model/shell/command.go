package shell

import (
	"bytes"
	"fmt"

	"github.com/jimschubert/docked/model/docker"
	"github.com/jimschubert/docked/model/docker/commands"
	d "github.com/moby/buildkit/frontend/dockerfile/parser"
	"mvdan.cc/sh/v3/syntax"
)

// PosixCommand is a simple representation of a command - the name of the command and any args passed to it
type PosixCommand struct {
	Name string
	Args []string
}

// NewPosixCommandFromNode extracts the "command" part of a Docker instruction.
func NewPosixCommandFromNode(node *d.Node) ([]PosixCommand, error) {
	dockerCommand, commandText := docker.Instruction(node)
	if dockerCommand != commands.Run {
		return []PosixCommand{}, fmt.Errorf("unexpected docker command: %v", dockerCommand)
	}
	return NewPosixCommand(commandText)
}

// NewPosixCommand parses input into an array of commands represented within that input
func NewPosixCommand(input string) ([]PosixCommand, error) {
	parser := syntax.NewParser(syntax.KeepComments(true))
	syntax.Variant(syntax.LangPOSIX)(parser)
	parsed, err := parser.Parse(bytes.NewReader([]byte(input)), "")
	if err != nil {
		return nil, err
	}

	commands := make([]PosixCommand, 0)
	syntax.Walk(parsed, func(node syntax.Node) bool {
		switch t := node.(type) {
		case *syntax.CallExpr:
			args := make([]string, 0)
			command := PosixCommand{}
			for i, arg := range t.Args {
				if i == 0 {
					command.Name = arg.Lit()
				} else {
					args = append(args, arg.Lit())
				}
			}
			command.Args = args

			commands = append(commands, command)
		default:
		}
		return true
	})
	return commands, nil
}
