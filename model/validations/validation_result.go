package validations

import (
	"github.com/jimschubert/docked/model"
)

type ValidationResult struct {
	Priority model.Priority
	Result   model.Valid
	Details  string
}

func NewFailureResult(name string, priority model.Priority, details string) *ValidationResult {
	return &ValidationResult{
		Priority: priority,
		Result:   model.Failure,
		Details:  details,
	}
}

func NewSuccessResult(name string, priority model.Priority, details string) *ValidationResult {
	return &ValidationResult{
		Priority: priority,
		Result:   model.Success,
		Details:  details,
	}
}
