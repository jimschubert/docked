package docked

import (
	"io/ioutil"
	"os"

	"github.com/jimschubert/docked/model"
	"gopkg.in/yaml.v3"
)

type RuleOverrides []ConfigRuleOverride

type Config struct {
	Ignore []string `yaml:"ignore"`
	// todo: support key: value as well
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

type ConfigRuleOverride struct {
	ID       string          `yaml:"id"`
	Priority *model.Priority `yaml:"priority,omitempty"`
}

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
