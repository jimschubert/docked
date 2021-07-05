package model

//go:generate stringer -type=Valid
// Valid is represents an enum of valid states for a validations.ValidationResult
// State names and uses are self-explanatory:
//	* Success
// 	* Failure
//	* Ignored
//	* Skipped
type Valid int

const (
	Success Valid = iota
	Failure
	Ignored
	Skipped
)
