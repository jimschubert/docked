package model

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Priority represents an enum of supported priorities.
// Priorities are more or less self-explanatory:
//  * LowPriority
//  * MediumPriority
//  * HighPriority
//  * CriticalPriority
//go:generate stringer -type=Priority
type Priority int

//goland:noinspection ALL
const (
	// LowPriority indicates less important rules, such as recommendations.
	LowPriority Priority = iota
	// MediumPriority indicates a rule which should be handled eventually, or which may be irrelevant based on
	// the base image or other factors which can't be evaluated statically. This can often be considered a "notice"
	// or cleanup task item.
	MediumPriority
	// HighPriority indicates a rule which raises non-security concerns.
	// Deployments with failed HighPriority rules can be done with good reason.
	HighPriority
	// CriticalPriority indicates a rule which raises a potential security or "correctness" conern which should be fixed
	// before deploying an image based on the current Dockerfile. Deploying an image with CriticalPriority issues could
	// result in production bugs or breaking consumers of the image.
	CriticalPriority
)

// Ptr is a utility function to return a pointer to the Priority pointer
func (i Priority) Ptr() *Priority {
	return &i
}

func (i *Priority) unmarshal(bytes []byte) error {
	if len(bytes) < 3 {
		*i = LowPriority
		return nil
	}

	var result strings.Builder
	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') {
			result.WriteByte(b)
		}
	}

	original := string(bytes)
	lower := strings.ToLower(result.String())
	switch strings.TrimSuffix(lower, "priority") {
	case "low":
		*i = LowPriority
	case "medium":
		*i = MediumPriority
	case "high":
		*i = HighPriority
	case "critical":
		*i = CriticalPriority
	default:
		return fmt.Errorf("unrecognized priority %q", original)
	}
	return nil
}

// UnmarshalYAML implements the yaml.v3 interface for unmarshalling YAML
func (i *Priority) UnmarshalYAML(value *yaml.Node) error {
	var original string
	if err := value.Decode(&original); err != nil {
		return err
	}

	return i.unmarshal([]byte(original))
}
