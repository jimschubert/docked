package rules

import (
	"sync"

	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

// RuleList is a collection of commands.DockerCommand(s) and their rules.
type RuleList map[commands.DockerCommand]*[]validations.Rule

// AddRule associates the provided Rule with this RuleList.
func (r RuleList) AddRule(rule validations.Rule) {
	// Rules are added to each command they're interested in.
	// For example, if a rule needs to evaluate COPY and USER,
	// the same instance is applied to both of these keys
	for _, dockerCommand := range rule.GetCommands() {
		if rules, ok := r[dockerCommand]; ok {
			*rules = append(*rules, rule)
		} else {
			rulesList := []validations.Rule{rule}
			r[dockerCommand] = &rulesList
		}
	}
}

var defaultRuleList = RuleList{}
var lock = sync.Mutex{}

// AddRule registers a rule with the default list of rules.
func AddRule(rule validations.Rule) {
	lock.Lock()
	defer lock.Unlock()
	defaultRuleList.AddRule(rule)
}

// DefaultRules returns the underlying collect of default rules.
func DefaultRules() RuleList {
	lock.Lock()
	defer lock.Unlock()
	return defaultRuleList
}
