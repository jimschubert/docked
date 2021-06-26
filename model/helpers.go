package model

import (
	"strings"

	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func StringPtr(s string) *string {
	return &s
}

func IsBuilderFrom(node *parser.Node) bool {
	return node.Value == string(commands.From) && strings.Contains(node.Original, " as ")
}

func IsFinalFrom(node *parser.Node) bool {
	return node.Value == string(commands.From) && !IsBuilderFrom(node)
}
