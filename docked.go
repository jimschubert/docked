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
	IgnoreRules              []string
	SuppressBuildKitWarnings bool
}

func (d *Docked) Analyze(location string) ([]validations.Validation, error) {
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

	ignoreLookup := make(map[string]bool, 0)
	for _, ignore := range d.IgnoreRules {
		ignoreLookup[ignore] = true
	}

	activeRules := getActiveRules(ignoreLookup)
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
					Rule:             &rule,
				})
			}
		}
	}

	return validationsRan, nil
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
					Rule:             &rule,
				})
			} else {
				log.Debugf("Skipped %s at %s: %s", ruleId, locations, result.Details)
			}
		}
	}
}

func getActiveRules(ignoreLookup map[string]bool) rules.RuleList {
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
