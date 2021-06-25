package validations

import (
	"github.com/jimschubert/docked/model/docker"
)

type Validation struct {
	ID   string
	Path string
	Rule *Rule
	ValidationResult
}

type ValidationContext struct {
	Line      string
	Locations []docker.Location
}
