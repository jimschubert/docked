package validations

import (
	"bytes"

	"github.com/jimschubert/docked/model"
)

// ValidationResult is the primitive defaulting the result, details,
// and contextual information about a rule evaluation
type ValidationResult struct {
	Result   model.Valid         `json:"result,omitempty"`
	Details  string              `json:"details,omitempty"`
	Contexts []ValidationContext `json:"contexts,omitempty"`
}

func (v ValidationResult) GoString() string {
	buf := bytes.Buffer{}

	buf.WriteString(v.Result.String())
	if len(v.Contexts) > 0 {
		buf.WriteString(": [")
	}
	for _, context := range v.Contexts {
		if context.HasRecommendations {
			buf.WriteString("?")
		}
		if context.CausedFailure {
			buf.WriteString("!")
		}
		for i, location := range context.Locations {
			if i != 0 {
				buf.WriteString(":")
			}
			buf.WriteString(location.String())
		}
	}
	if len(v.Contexts) > 0 {
		buf.WriteString("]")
	}

	return buf.String()
}

// NewValidationResultSkipped is a utility function to create a result which is model.Skipped
func NewValidationResultSkipped(details string) *ValidationResult {
	return &ValidationResult{Result: model.Skipped, Details: details}
}

// NewValidationResultIgnored is a utility function to create a result which is model.Ignored
func NewValidationResultIgnored(details string) *ValidationResult {
	return &ValidationResult{Result: model.Ignored, Details: details}
}
