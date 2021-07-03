//go:generate go run ./cmd/gen.go
package docked

import (
	"os"
	"path/filepath"
	"sort"

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
	rulePriorityOverrides    *map[string]model.Priority
}

type AnalysisResult struct {
	Evaluated    []validations.Validation
	NotEvaluated []validations.Validation
}

type ConfiguredRules struct {
	Active   rules.RuleList
	Inactive rules.RuleList
}

func (d *Docked) AnalyzeWithRuleList(location string, configuredRules ConfiguredRules) (AnalysisResult, error) {
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
	validationsNotRan := make([]validations.Validation, 0)
	deferredEvaluationRules := make(map[string]validations.FinalizingRule)

	if !d.SuppressBuildKitWarnings {
		// This dumps out any warnings directly from buildkit to stdout
		//goland:noinspection GoNilness
		p.PrintWarnings(log.StandardLogger().Out)
	}

	seenCommands := make(map[commands.DockerCommand]bool, 0)

	//goland:noinspection ALL
	for _, node := range p.AST.Children {
		thisCommand := commands.DockerCommand(node.Value)
		seenCommands[thisCommand] = true
		if commandRules, ok := configuredRules.Active[thisCommand]; ok {
			if commandRules == nil {
				log.Warnf("Active rule mapped 0 rules to command %s", thisCommand)
				continue
			}
			currentRules := *commandRules
			d.evaluateNode(node, &currentRules, &validationsRan, &validationsNotRan, &deferredEvaluationRules, fullPath)
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

	for command, commandRules := range configuredRules.Active {
		if !seenCommands[command] {
			if commandRules != nil {
				for _, rule := range *commandRules {
					validationsNotRan = append(validationsNotRan, validations.Validation{
						ID:               rule.GetLintID(),
						Path:             fullPath,
						ValidationResult: *validations.NewValidationResultSkipped("The rule was not applicable to this Dockerfile"),
						Rule:             d.ruleCopy(rule),
					})
				}
			}
		}
	}

	for command, commandRules := range configuredRules.Inactive {
		if !seenCommands[command] {
			if commandRules != nil {
				for _, rule := range *commandRules {
					validationsNotRan = append(validationsNotRan, validations.Validation{
						ID:               rule.GetLintID(),
						Path:             fullPath,
						ValidationResult: *validations.NewValidationResultIgnored("The rule was ignored via configuration"),
						Rule:             d.ruleCopy(rule),
					})
				}
			}
		}
	}

	// Ensure returned lists are in consistent orders
	sort.Slice(validationsRan, func(left, right int) bool {
		return validationsRan[left].ID < validationsRan[right].ID
	})
	sort.Slice(validationsNotRan, func(left, right int) bool {
		return validationsNotRan[left].ID < validationsNotRan[right].ID
	})

	return AnalysisResult{Evaluated: validationsRan, NotEvaluated: validationsNotRan}, nil
}

func (d *Docked) Analyze(location string) (AnalysisResult, error) {
	configuredRules := buildConfiguredRules(d.Config)
	return d.AnalyzeWithRuleList(location, configuredRules)
}

func (d *Docked) evaluateNode(
	node *parser.Node,
	commandRules *[]validations.Rule,
	validationsRan *[]validations.Validation,
	validationsNotRan *[]validations.Validation,
	deferredRules *map[string]validations.FinalizingRule,
	fullPath string,
) {
	evaluating := *commandRules
	for _, rule := range evaluating {
		ruleId := rule.GetLintID()
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
				*validationsNotRan = append(*validationsNotRan, validations.Validation{
					ID:               ruleId,
					Path:             fullPath,
					ValidationResult: *result,
					Rule:             d.ruleCopy(rule),
				})
			}
		}
	}
}

func buildConfiguredRules(config Config) ConfiguredRules {
	ignoreLookup := make(map[string]bool, 0)
	for _, ignore := range config.Ignore {
		ignoreLookup[ignore] = true
	}

	// ConfiguredRules
	activeRules := rules.RuleList{}
	inactiveRules := rules.RuleList{}
	for _, r := range rules.DefaultRules() {
		for _, rule := range *r {
			ruleId := rule.GetLintID()
			if ignoreLookup[rule.GetLintID()] {
				log.Debugf("Ignoring rule %s", ruleId)
				inactiveRules.AddRule(rule)
			} else {
				if resettable, ok := rule.(validations.ResettingRule); ok {
					resettable.Reset()
				}
				activeRules.AddRule(rule)
			}
		}
	}
	return ConfiguredRules{Active: activeRules, Inactive: inactiveRules}
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

	priority := r.GetPriority()
	if override, ok := (*d.rulePriorityOverrides)[r.GetLintID()]; ok {
		priority = override
	}

	var rule validations.Rule = validations.SimpleRule{
		Name:     r.GetName(),
		Summary:  r.GetSummary(),
		Details:  r.GetDetails(),
		Priority: priority,
		Commands: r.GetCommands(),
		Handler: func(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
			log.Warnf("Rule %s is only intended to be invoked via the Analyze function", r.GetLintID())
			return nil
		},
		Category: r.GetCategory(),
		URL:      r.GetURL(),
	}
	return &rule
}
