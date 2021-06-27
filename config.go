package docked

import "github.com/jimschubert/docked/model"

type RuleOverrides []ConfigRuleOverride

type Config struct {
	Ignore        []string              `yaml:"ignore"`
	// todo: support key: value as well
	RuleOverrides *RuleOverrides `yaml:"rule_overrides,omitempty"`
}

type ConfigRuleOverride struct {
	ID       string          `yaml:"id"`
	Priority *model.Priority `yaml:"priority,omitempty"`
}
