package model

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:generate stringer -type=Priority
type Priority int

//goland:noinspection ALL
const (
	LowPriority Priority = iota
	MediumPriority
	HighPriority
	CriticalPriority
)

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

func (i *Priority) UnmarshalYAML(value *yaml.Node) error {
	var original string
	if err := value.Decode(&original); err != nil {
		return err
	}

	return i.unmarshal([]byte(original))
}
