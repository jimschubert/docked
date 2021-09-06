package docked

import (
	"io/ioutil"
	"os"
	"sort"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/validations"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// ConfigRuleOverride defines the id-priority override mapping used in a config file
type ConfigRuleOverride struct {
	// The rule id to override
	ID string `yaml:"id"`
	// The overridden priority
	Priority *model.Priority `yaml:"priority,omitempty"`
}

// RuleOverrides is a slice of ConfigRuleOverride. This type allows for simpler definitions and YAML parsing.
type RuleOverrides []ConfigRuleOverride

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

// Config represents the YAML config structure exposed to users
type Config struct {
	// Ignore this collection of rule ids
	Ignore []string `yaml:"ignore"`
	// RuleOverrides allows users to override the ConfigRuleOverride.Priority of a specific rule by ConfigRuleOverride.ID
	RuleOverrides *RuleOverrides `yaml:"rule_overrides,omitempty"`
	CustomRules []validations.SimpleRegexRule `yaml:"custom_rules,omitempty"`
}

// Load a Config from path with sorted members (by ID for rule overrides, by Name for custom rules)
func (c *Config) Load(path string) error {
	b, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		log.Warnf("Config file does not exist! Attempted: %s", path)
		return nil
	}

	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(b, c); err != nil {
		return err
	}

	// Sorting each slice in config. Necessary because yaml v3 doesn't always return top-down parsing order.
	if len(c.Ignore) > 0 {
		sort.Strings(c.Ignore)
	}

	if c.RuleOverrides != nil {
		overrides := *c.RuleOverrides
		sort.Slice(overrides, func(i, j int) bool {
			return overrides[i].ID < overrides[j].ID
		})
		c.RuleOverrides = &overrides
	}

	if len(c.CustomRules) > 0 {
		sort.Slice(c.CustomRules, func(i, j int) bool {
			return c.CustomRules[i].Name < c.CustomRules[j].Name
		})
	}

	return nil
}
