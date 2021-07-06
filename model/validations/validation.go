package validations

import (
	"github.com/jimschubert/docked/model/docker"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Validation struct {
	ID               string `json:"id,omitempty"`
	Path             string `json:"path,omitempty"`
	Rule             *Rule  `json:"rule,omitempty"`
	ValidationResult `json:"validation_result"`
}

type ValidationContext struct {
	Line          string            `json:"line,omitempty"`
	Locations     []docker.Location `json:"locations,omitempty"`
	CausedFailure bool              `json:"caused_failure,omitempty"`
}

type NodeValidationContext struct {
	Node    parser.Node
	Context ValidationContext
}
