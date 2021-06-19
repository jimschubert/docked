package rules

import (
	"sync"

	"github.com/jimschubert/docked/model/docker/command"
)

type RuleList map[command.DockerCommand]*[]Rule

var defaultRuleList = RuleList{}
var lock = sync.Mutex{}

func AddRule(rule Rule) {
	lock.Lock()
	defer lock.Unlock()
	if rules, ok := defaultRuleList[rule.Command]; ok {
		*rules = append(*rules, rule)
	} else {
		rulesList := []Rule{rule}
		defaultRuleList[rule.Command] = &rulesList
	}
}

func DefaultRules() RuleList {
	lock.Lock()
	defer lock.Unlock()
	return defaultRuleList
}
