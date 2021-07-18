package docked

import (
	"io/ioutil"
	"os"

	"github.com/jimschubert/docked/model"
	"gopkg.in/yaml.v3"
)

// RuleOverrides is a slice of ConfigRuleOverride. This type allows for simpler definitions and YAML parsing.
type RuleOverrides []ConfigRuleOverride

// Config represents the YAML config structure exposed to users
type Config struct {
	// Ignore this collection of rule ids
	Ignore []string `yaml:"ignore"`
	// RuleOverrides allows users to override the ConfigRuleOverride.Priority of a specific rule by ConfigRuleOverride.ID
	RuleOverrides *RuleOverrides `yaml:"rule_overrides,omitempty"`
}

// Load a Config from path
func (c *Config) Load(path string) error {
	b, err := ioutil.ReadFile(path)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, c)
}

// ConfigRuleOverride defines the id-priority override mapping used in a config file
type ConfigRuleOverride struct {
	// The rule id to override
	ID string `yaml:"id"`
	// The overridden priority
	Priority *model.Priority `yaml:"priority,omitempty"`
}

// UnmarshalYAML implements the interface necessary to have greater control over deserializing RuleOverrides
func (r *RuleOverrides) UnmarshalYAML(value *yaml.Node) error {
	*r = make([]ConfigRuleOverride, 0)
	var kvp map[string]model.Priority
	if err := value.Decode(&kvp); err == nil {
		for key, priority := range kvp {
			*r = append(*r, ConfigRuleOverride{key, priority.Ptr()})
		}
		return nil
	}

	type raw RuleOverrides
	return value.Decode((*raw)(r))
}
