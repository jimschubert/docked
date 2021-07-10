package model

import (
	"strings"

	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// StringPtr obtains a pointer of provided string input
func StringPtr(s string) *string {
	return &s
}

// IsBuilderFrom determines if the current node is a FROM command defining an explicit builder syntax
// Note that this doesn't account for index-based builder patterns.
func IsBuilderFrom(node *parser.Node) bool {
	return node.Value == string(commands.From) && strings.Contains(node.Original, " as ")
}

// IsFinalFrom determines whether the current node is a FROM command without the explicit alias used in named builder pattern
func IsFinalFrom(node *parser.Node) bool {
	return node.Value == string(commands.From) && !IsBuilderFrom(node)
}
