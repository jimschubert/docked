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
	deferredEvaluationRules := make([]rules.Rule, 0)

	if !d.SuppressBuildKitWarnings {
		// This dumps out any warnings directly from buildkit to stdout
		//goland:noinspection GoNilness
		p.PrintWarnings(log.StandardLogger().Out)
	}

	//goland:noinspection ALL
	for _, node := range p.AST.Children {
		thisCommand := commands.DockerCommand(node.Value)
		if commandRules, ok := activeRules[thisCommand]; ok {
			d.evaluateNode(node, commandRules, &validationsRan, &deferredEvaluationRules, fullPath)
		}
	}

	if len(deferredEvaluationRules) > 0 {
		for _, rule := range deferredEvaluationRules {
			result := rule.InvokeFinalize()
			validationsRan = append(validationsRan, validations.Validation{
				ID:               rule.LintID(),
				Path:             fullPath,
				ValidationResult: *result,
			})
		}
	}

	return validationsRan, nil
}

func (d *Docked) evaluateNode(
	node *parser.Node,
	commandRules *[]rules.Rule,
	validationsRan *[]validations.Validation,
	deferredRules *[]rules.Rule,
	fullPath string,
) {
	for _, rule := range *commandRules {
		ruleId := rule.LintID()
		locations := docker.FromParserRanges(node.Location())
		validationContext := validations.ValidationContext{
			Line:      node.Original,
			Locations: locations,
		}

		result := rule.Evaluate(node, validationContext)
		if rule.HasFinalizer() {
			*deferredRules = append(*deferredRules, rule)
			return
		}

		if result.Result != model.Skipped {
			*validationsRan = append(*validationsRan, validations.Validation{
				ID:               ruleId,
				Path:             fullPath,
				ValidationResult: *result,
			})
		} else {
			log.Debugf("Skipped %s at %s: %s", ruleId, locations, result.Details)
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
				rule.InvokeReset()
				activeRules.AddRule(rule)
			}
		}
	}
	return activeRules
}
