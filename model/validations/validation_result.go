package validations

import (
	"github.com/jimschubert/docked/model"
)

type ValidationResult struct {
	Result   model.Valid         `json:"result,omitempty"`
	Details  string              `json:"details,omitempty"`
	Contexts []ValidationContext `json:"contexts,omitempty"`
}

func NewValidationResultSkipped(details string) *ValidationResult {
	return &ValidationResult{Result: model.Skipped, Details: details}
}

func NewValidationResultIgnored(details string) *ValidationResult {
	return &ValidationResult{Result: model.Ignored, Details: details}
}
