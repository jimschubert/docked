package docked

//go:generate go run ./cmd/generators/rules_md.go
import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/rules"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

// Docked is the main type for initializing Dockerfile linting/analysis
type Docked struct {
	// Configuration for analysis
	Config Config
	// Suppress the underlying warnings presented by buildkit's parser. Use this if you want to pipe text summary to file.
	SuppressBuildKitWarnings bool
	rulePriorityOverrides    *map[string]model.Priority
}

// AnalysisResult holds final validations, separated in those which have been Evaluated and those which have not (NotEvaluated).
// A validations.Validation holds references to the rule and the result of validation to simplify reporting.
type AnalysisResult struct {
	Evaluated    []validations.Validation `json:"evaluated"`
	NotEvaluated []validations.Validation `json:"not_evaluated"`
}

// GoString returns a string representation for formatter patterns %#v
func (a AnalysisResult) GoString() string {
	buf := bytes.Buffer{}
	if len(a.Evaluated) > 0 {
		buf.WriteString("E(")
		for i, validation := range a.Evaluated {
			if i != 0 {
				buf.WriteString("|")
			}
			buf.WriteString(fmt.Sprintf("%#v", validation))
		}
		buf.WriteString(")")
	}
	if len(a.Evaluated) > 0 && len(a.NotEvaluated) > 0 {
		buf.WriteString(" ")
	}
	if len(a.NotEvaluated) > 0 {
		buf.WriteString("N(")
		for i, validation := range a.NotEvaluated {
			if i != 0 {
				buf.WriteString("|")
			}
			buf.WriteString(fmt.Sprintf("%#v", validation))
		}
		buf.WriteString(")")
	}
	return buf.String()
}

// ConfiguredRules partitions results into active and inactive lists
type ConfiguredRules struct {
	Active   rules.RuleList
	Inactive rules.RuleList
}

// AnalyzeWithRuleList is just like Analyze, but accepts an additional parameter of ConfiguredRules
//
// This allows programmatic evaluation of rules without the ignore/priority overrides done as a default within Analyze.
//
// Returns the AnalysisResult or error.
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

	seenCommands := make(map[commands.DockerCommand]bool)

	finalStage := d.finalStageIndex(p.AST.Children)

	//goland:noinspection ALL
	for idx, node := range p.AST.Children {
		thisCommand := commands.Of(node.Value)
		seenCommands[thisCommand] = true
		if commandRules, ok := configuredRules.Active[thisCommand]; ok {
			if commandRules == nil {
				log.Warnf("Active rule mapped 0 rules to command %s", thisCommand)
				continue
			}
			currentRules := *commandRules
			isBuilderStage := idx < finalStage
			d.evaluateNode(node, isBuilderStage, &currentRules, &validationsRan, &validationsNotRan, &deferredEvaluationRules, fullPath)
		}
	}

	if len(deferredEvaluationRules) > 0 {
		for ruleID, finalizer := range deferredEvaluationRules {
			log.Tracef("Evaluating deferred rule %s", ruleID)
			result := finalizer.Finalize()
			if result != nil {
				rule := finalizer.(validations.Rule)
				validationsRan = append(validationsRan, validations.Validation{
					ID:               ruleID,
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

// Analyze a dockerfile residing at location.
//
// All known rules which are applicable to the Dockerfile contents are evaluated,
// allowing configuration-based ignores and manipulation of priority/severity of rules.
//
// Returns the AnalysisResult or error.
func (d *Docked) Analyze(location string) (AnalysisResult, error) {
	configuredRules := buildConfiguredRules(d.Config)
	return d.AnalyzeWithRuleList(location, configuredRules)
}

// finalStageIndex is a preprocessor which evaluates nodes in reverse to determine where the final build context
// starts (last index of FROM). This allows evaluation to also handle index-based builder contexts for rules where AppliesToBuilder is false.
func (d Docked) finalStageIndex(nodes []*parser.Node) int {
	var finalStageAt int
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		if commands.Of(node.Value) == commands.From {
			finalStageAt = i
			break
		}
	}
	return finalStageAt
}

// evaluateNode invokes rule evaluation. It determines whether the evaluated rule should be deferred, and partitions into ran/notRan collections.
func (d *Docked) evaluateNode(
	node *parser.Node,
	isBuilderStage bool,
	commandRules *[]validations.Rule,
	validationsRan *[]validations.Validation,
	validationsNotRan *[]validations.Validation,
	deferredRules *map[string]validations.FinalizingRule,
	fullPath string,
) {
	evaluating := *commandRules
	for _, rule := range evaluating {
		ruleID := rule.GetLintID()
		locations := docker.FromParserRanges(node.Location())
		validationContext := validations.ValidationContext{
			Line:             node.Original,
			Locations:        locations,
			IsBuilderContext: isBuilderStage,
		}

		result := rule.Evaluate(node, validationContext)
		if finalizer, ok := rule.(validations.FinalizingRule); ok {
			// add the rule as deferred if we haven't yet seen it
			if _, ok := (*deferredRules)[ruleID]; !ok {
				log.Tracef("Deferring evaluation of rule %s at line %d", ruleID, validationContext.Locations[0].Start.Line)
				(*deferredRules)[ruleID] = finalizer
			}
			continue
		}

		if result != nil {
			if result.Result != model.Skipped {
				v := validations.Validation{
					ID:               ruleID,
					Path:             fullPath,
					ValidationResult: *result,
					Rule:             d.ruleCopy(rule),
				}
				printValidationResults(v)
				*validationsRan = append(*validationsRan, v)
			} else {
				log.Tracef("Skipped %s at %s: %s", ruleID, locations, result.Details)
				v := validations.Validation{
					ID:               ruleID,
					Path:             fullPath,
					ValidationResult: *result,
					Rule:             d.ruleCopy(rule),
				}
				printRulesSkipped(v)
				*validationsNotRan = append(*validationsNotRan, v)
			}
		}
	}
}

// printValidationResults formats a debug message for a processed validation/rule
func printValidationResults(v validations.Validation) {
	var indicator string
	r := *v.Rule
	priority := strings.TrimSuffix(r.GetPriority().String(), "Priority")
	if v.ValidationResult.Result == model.Success {
		indicator = "✔"
		var lineInfo = ""
		if len(v.Contexts) > 0 {
			lineInfo = fmt.Sprintf("\n\t%s> %s", v.Contexts[0].Locations, v.Contexts[0].Line)
		}
		log.Debugf("%s %-8s %s %s\n\t%s", indicator, priority, v.ID, lineInfo, v.Details)
	} else {
		indicator = "⨯"
		var where validations.ValidationContext
		// grab the first hit. Other reporting will reference all locations with issues.
		for _, context := range v.Contexts {
			if context.CausedFailure {
				where = context
				break
			}
		}
		log.Debugf("%s %-8s %s \n\t%s> %s\n\t%s", indicator, priority, v.ID, where.Locations, where.Line, v.Details)
	}
}

// printRulesSkipped formats a debug message for a skipped validation/rule
func printRulesSkipped(v validations.Validation) {
	indicator := "#"
	r := *v.Rule
	priority := strings.TrimSuffix(r.GetPriority().String(), "Priority")
	log.Debugf("%s %-8s %s \n\t%s", indicator, priority, v.ID, v.Details)
}

// buildConfiguredRules evaluates which rules to ignore via config, and splits all known rules into active and inactive collections, exposed as ConfiguredRules
func buildConfiguredRules(config Config) ConfiguredRules {
	ignoreLookup := make(map[string]bool)
	includeLookup := make(map[string]bool)
	for _, ignore := range config.Ignore {
		ignoreLookup[ignore] = true
	}
	for _, include := range config.IncludeRules {
		includeLookup[include] = true
	}

	// ConfiguredRules
	activeRules := rules.RuleList{}
	inactiveRules := rules.RuleList{}

	for _, r := range rules.DefaultRules() {
		for _, rule := range *r {
			ruleID := rule.GetLintID()
			if ignoreLookup[ruleID] {
				log.Debugf("Ignoring rule %s", ruleID)
				inactiveRules.AddRule(rule)
				// only need to account for multi-command rules once
				ignoreLookup[ruleID] = false
			} else {
				if resettable, ok := rule.(validations.ResettingRule); ok {
					resettable.Reset()
				}
				if !config.SkipDefaultRules {
					activeRules.AddRule(rule)
				} else if includeLookup[ruleID] {
					activeRules.AddRule(rule)
					// only need to account for multi-command rules once
					includeLookup[ruleID] = false
				}
			}
		}
	}

	for _, customRule := range config.CustomRules {
		activeRules.AddRule(customRule)
	}

	return ConfiguredRules{Active: activeRules, Inactive: inactiveRules}
}

// ruleCopy allows to expose a Rule to caller of Analyze without exposing its handler to be re-run.
// The caller is still allowed to invoke Evaluate from default rules. This copy is intended only to communicate
// the expectation that rule evaluation occurs through Analyze or other working directly on the rule list.
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
		log.Debugf("Overriding %s priority to %s", r.GetLintID(), override.String())
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
