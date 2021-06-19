package model

//go:generate stringer -type=Priority
type Priority int

//goland:noinspection ALL
const (
	LowPriority Priority = iota
	MediumPriority
	HighPriority
	CriticalPriority
)
