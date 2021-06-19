package validations

import "github.com/jimschubert/docked/model/docker"

type Validation struct {
	ID string
	ValidationResult
	Line  string
	Range []docker.Range
}
