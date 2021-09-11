package validations

import (
	"github.com/jimschubert/docked/model/docker"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Validation represents the result of evaluating the Rule against the Dockerfile located at Path
type Validation struct {
	ID               string                     `json:"id,omitempty"`   // The ID of the rule
	Path             string                     `json:"path,omitempty"` // The Path of the Dockerfile
	Rule             *Rule                      `json:"rule,omitempty"` // The Rule applied
	ValidationResult `json:"validation_result"` // A Validation is composed of ValidationResult
}

// ValidationContext details whether a line and positions within that line caused a validation failure
type ValidationContext struct {
	Line               string            `json:"line,omitempty"`                // The Line of text being evaluated, as parsed from the Dockerfile
	Locations          []docker.Location `json:"locations,omitempty"`           // The start and end Locations within the Dockerfile
	CausedFailure      bool              `json:"caused_failure,omitempty"`      // Whether the parsed Line caused a failure in the final Validation
	HasRecommendations bool              `json:"has_recommendations,omitempty"` // Whether the parsed Line includes a recommendation in the final Validation
}

// NodeValidationContext associates a parser.Node and ValidationContext, such as deferred execution via rules implementing FinalizingRule.
type NodeValidationContext struct {
	Node    parser.Node
	Context ValidationContext
}
