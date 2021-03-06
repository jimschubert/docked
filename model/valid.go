package model

// Valid is represents an enum of valid states for a validations.ValidationResult
// State names and uses are self-explanatory:
//	* Success
// 	* Failure
//	* Ignored
//	* Skipped
//  * Recommendation
//go:generate stringer -type=Valid
type Valid int

const (
	// Success is when validation results in a non-failure state.
	// This does not imply the validation was relevant to the line(s) evaluated, which is up to the rule.
	Success Valid = iota
	// Failure is when validation was not successful.
	Failure
	// Ignored is when a rule is not evaluated because it was specifically ignored by the user.
	Ignored
	// Skipped is when a rule is not evaluated, but was not skipped by the user (in which case it will be Ignored).
	// A Skipped state may occur when a rule is only contextually relevant, such as being irrelevant in builder contexts.
	Skipped

	// Recommendation is a completed validation that is neither success/failure but more of a best practice.
	// A result returning Recommendation should not cause any failures on execution. These represent situations which are
	// generally not a good idea (e.g. ADD http://example.com/root.tgz) but have valid use cases (e.g. scratch images).
	Recommendation
)
