package validations

import (
	"github.com/jimschubert/docked/model"
)

// ValidationResult is the primitive defaulting the result, details,
// and contextual information about a rule evaluation
type ValidationResult struct {
	Result   model.Valid         `json:"result,omitempty"`
	Details  string              `json:"details,omitempty"`
	Contexts []ValidationContext `json:"contexts,omitempty"`
}

// NewValidationResultSkipped is a utility function to create a result which is model.Skipped
func NewValidationResultSkipped(details string) *ValidationResult {
	return &ValidationResult{Result: model.Skipped, Details: details}
}

// NewValidationResultIgnored is a utility function to create a result which is model.Ignored
func NewValidationResultIgnored(details string) *ValidationResult {
	return &ValidationResult{Result: model.Ignored, Details: details}
}
