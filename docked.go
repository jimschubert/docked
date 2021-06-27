//go:generate go run ./cmd/gen.go
package docked

import (
	"os"
	"path/filepath"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/rules"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

type Docked struct {
	Config                   Config
	SuppressBuildKitWarnings bool
	rulePriorityOverrides *map[string]model.Priority
}

func (d *Docked) AnalyzeWithRuleList(location string, activeRules rules.RuleList) ([]validations.Validation, error) {
	var err error
	fullPath, err := filepath.Abs(location)
	if err != nil {
		log.Fatal("Could not determine absolute path to Dockerfile")
	}

	dockerfile, err := os.Open(fullPath)
	if err != nil {
		log.Fatal("Could not open path")
	}
	p, err := parser.Parse(dockerfile)
	if err != nil || p == nil {
		log.Fatal("Could not parse Dockerfile")
	}

	validationsRan := make([]validations.Validation, 0)
	deferredEvaluationRules := make(map[string]validations.FinalizingRule)

	if !d.SuppressBuildKitWarnings {
		// This dumps out any warnings directly from buildkit to stdout
		//goland:noinspection GoNilness
		p.PrintWarnings(log.StandardLogger().Out)
	}

	//goland:noinspection ALL
	for _, node := range p.AST.Children {
		thisCommand := commands.DockerCommand(node.Value)
		if commandRules, ok := activeRules[thisCommand]; ok {
			currentRules := *commandRules
			d.evaluateNode(node, &currentRules, &validationsRan, &deferredEvaluationRules, fullPath)
		}
	}

	if len(deferredEvaluationRules) > 0 {
		for lintId, finalizer := range deferredEvaluationRules {
			result := finalizer.Finalize()
			if result != nil {
				rule := finalizer.(validations.Rule)
				validationsRan = append(validationsRan, validations.Validation{
					ID:               lintId,
					Path:             fullPath,
					ValidationResult: *result,
					Rule:             d.ruleCopy(rule),
				})
			}
		}
	}

	return validationsRan, nil
}

func (d *Docked) Analyze(location string) ([]validations.Validation, error) {
	activeRules := getActiveRules(d.Config)
	return d.AnalyzeWithRuleList(location, activeRules)
}

func (d *Docked) evaluateNode(
	node *parser.Node,
	commandRules *[]validations.Rule,
	validationsRan *[]validations.Validation,
	deferredRules *map[string]validations.FinalizingRule,
	fullPath string,
) {
	evaluating := *commandRules
	for _, rule := range evaluating {
		ruleId := rule.LintID()
		locations := docker.FromParserRanges(node.Location())
		validationContext := validations.ValidationContext{
			Line:      node.Original,
			Locations: locations,
		}

		result := rule.Evaluate(node, validationContext)
		if finalizer, ok := rule.(validations.FinalizingRule); ok {
			if _, ok := (*deferredRules)[ruleId]; !ok {
				(*deferredRules)[ruleId] = finalizer
			}
			continue
		}

		if result != nil {
			if result.Result != model.Skipped {
				*validationsRan = append(*validationsRan, validations.Validation{
					ID:               ruleId,
					Path:             fullPath,
					ValidationResult: *result,
					Rule:             d.ruleCopy(rule),
				})
			} else {
				log.Debugf("Skipped %s at %s: %s", ruleId, locations, result.Details)
			}
		}
	}
}

func getActiveRules(config Config) rules.RuleList {
	ignoreLookup := make(map[string]bool, 0)
	for _, ignore := range config.Ignore {
		ignoreLookup[ignore] = true
	}

	activeRules := rules.RuleList{}
	for _, r := range rules.DefaultRules() {
		for _, rule := range *r {
			ruleId := rule.LintID()
			if ignoreLookup[rule.LintID()] {
				log.Debugf("Ignoring rule %s", ruleId)
			} else {
				if resettable, ok := rule.(validations.ResettingRule); ok {
					resettable.Reset()
				}
				activeRules.AddRule(rule)
			}
		}
	}
	return activeRules
}

// ruleCopy allows to expose a Rule to caller of docked.Analyze without exposing its handler to be re-run.
// The caller is still allowed to invoke Evaluate from default rules. This copy is intended only to communicate
// the expectation that rule evaluation occurs through docked.Analyze or other working directly on the rule list.
func (d *Docked) ruleCopy(r validations.Rule) *validations.Rule {
	if d.rulePriorityOverrides == nil {
		overrides := make(map[string]model.Priority)
		if d.Config.RuleOverrides != nil {
			for _, override := range *d.Config.RuleOverrides {
				if override.Priority != nil {
					overrides[override.ID] = *override.Priority
				}
			}
		}
		(*d).rulePriorityOverrides = &overrides
	}

	priority := r.Priority()
	if override, ok := (*d.rulePriorityOverrides)[r.LintID()]; ok {
		priority = override
	}

	rule := validations.NewSimpleRule(
		r.Name(),
		r.Summary(),
		r.Details(),
		priority,
		r.Commands(),
		func(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
			log.Warnf("Rule %s is only intended to be invoked via the Analyze function", r.LintID())
			return nil
		},
		r.Category(),
		r.URL(),
	)
	return &rule
}
