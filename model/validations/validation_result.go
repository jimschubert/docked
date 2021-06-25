package validations

import (
	"github.com/jimschubert/docked/model"
)

type ValidationResult struct {
	Result   model.Valid
	Details  string
	Contexts []ValidationContext
}

func NewSkippedResult(details string) *ValidationResult {
	return &ValidationResult{Result: model.Skipped, Details: details}
}