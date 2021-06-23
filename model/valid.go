package model

//go:generate stringer -type=Valid
type Valid int

const (
	Success Valid = iota
	Failure
	Ignored
	Skipped
)
