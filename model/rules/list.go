package rules

import (
	"sync"

	"github.com/jimschubert/docked/model/docker/commands"
)

type RuleList map[commands.DockerCommand]*[]Rule

func (r RuleList) AddRule(rule Rule) {
	// Rules are added to each command they're interested in.
	// For example, if a rule needs to evaluate COPY and USER,
	// the same instance is applied to both of these keys
	for _, dockerCommand := range rule.Commands {
		if rules, ok := r[dockerCommand]; ok {
			*rules = append(*rules, rule)
		} else {
			rulesList := []Rule{rule}
			r[dockerCommand] = &rulesList
		}
	}
}

var defaultRuleList = RuleList{}
var lock = sync.Mutex{}

func AddRule(rule Rule) {
	lock.Lock()
	defer lock.Unlock()
	defaultRuleList.AddRule(rule)
}

func DefaultRules() RuleList {
	lock.Lock()
	defer lock.Unlock()
	return defaultRuleList
}
