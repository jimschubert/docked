package validations

import (
	"github.com/jimschubert/docked/model/docker"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Validation struct {
	ID   string
	Path string
	Rule *Rule
	ValidationResult
}

type ValidationContext struct {
	Line      string
	Locations []docker.Location
}

type NodeValidationContext struct {
	Node    parser.Node
	Context ValidationContext
}
