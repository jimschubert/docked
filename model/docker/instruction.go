package docker

import (
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Instruction extracts the commands.DockerCommand and instruction text from a Docker instruction
func Instruction(node *parser.Node) (commands.DockerCommand, string) {
	trimStart := len(node.Value) + 1 // command plus trailing space
	commandText := node.Original[trimStart:]

	return commands.Of(node.Value), commandText
}
