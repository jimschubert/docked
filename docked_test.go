package docked

import (
	"bytes"
	"fmt"
	"os"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/rules"
	"github.com/jimschubert/docked/model/validations"
)

func printEvaluated(evaluated []validations.Validation) {
	for _, r := range evaluated {
		b := bytes.Buffer{}
		_, _ = b.WriteString(fmt.Sprintf("%s - %s", r.ID, r.Result.String()))
		if r.Result == model.Failure {
			for _, c := range r.Contexts {
				if c.CausedFailure {
					_, _ = b.WriteString(fmt.Sprintf(" * [%2d] %s\n", c.Locations[0].Start.Line, c.Line))
				}
			}
		} else {
			_, _ = b.WriteString("\n")
		}
		_, _ = fmt.Fprint(os.Stdout, b.String())
	}
}

// ExampleDocked_Analyze provides an example of programmatically invoking Docked.Analyze with default rules
func ExampleDocked_Analyze() {
	c := Config{}
	if err := c.Load("./testdata/config/example.yml"); err != nil {
		panic("Failed to load config file!")
	}

	d := Docked{
		Config:                   c,
		SuppressBuildKitWarnings: true,
	}

	result, err := d.Analyze("./testdata/minimal.dockerfile")
	if err != nil {
		panic("Failed to analyze dockerfile")
	}

	// programmatically consume array of evaluated and/or not-evaluated rules
	printEvaluated(result.Evaluated)

	// Output:
	// D5:no-debian-frontend - Success
	// D5:secret-aws-access-key - Success
	// D5:secret-aws-secret-access-key - Success
	// DC:consider-multistage - Success
	// DC:curl-without-fail - Success
	// DC:gpg-without-batch - Success
	// DC:gpg-without-batch - Success
	// DC:layered-ownership-change - Success
}

// ExampleDocked_AnalyzeWithRuleList provides an example of programmatically invoking Docked.AnalyzeWithRuleList
// with user-defined rules. See also reporter.TextReporter and reporter.HTMLReporter for in-built output formatters.
func ExampleDocked_AnalyzeWithRuleList() {
	d := Docked{}

	// user can extend default rule set or define their own
	activeRules := rules.RuleList{}
	myRule := validations.SimpleRegexRule{
		Name:     "no-distroless",
		Pattern:  `\bgcr\.io/distroless\b`,
		Priority: model.CriticalPriority,
		Command:  commands.From,
	}
	activeRules.AddRule(myRule)

	result, err := d.AnalyzeWithRuleList("./testdata/minimal.dockerfile", ConfiguredRules{Active: activeRules})
	if err != nil {
		panic("Failed to analyze dockerfile")
	}

	// programmatically consume array of evaluated and/or not-evaluated rules
	printEvaluated(result.Evaluated)

	// Output:
	// D7:no-distroless - Success
	// D7:no-distroless - Failure * [13] FROM gcr.io/distroless/base-debian10
}
