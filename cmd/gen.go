// +build generate

package main

import (
	"html/template"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jimschubert/docked/model/rules"
	"github.com/jimschubert/docked/model/validations"
)

func main() {
	generateRulesReadme()
}

func generateRulesReadme() {
	f, err := os.Create("RULES.md")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	allRules := make([]validations.Rule, 0)
	seenRule := make(map[string]bool, 0)
	rulesList := rules.DefaultRules()

	for _, commandRules := range rulesList {
		if commandRules != nil {
			for _, rule := range *commandRules {
				id := rule.LintID()
				if !seenRule[id] {
					allRules = append(allRules, rule)
					seenRule[id] = true
				}
			}
		}
	}

	sort.Slice(allRules, func(i, j int) bool {
		return allRules[i].LintID() < allRules[j].LintID()
	})

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"shortPriority": func(s string) string {
			return strings.TrimSuffix(s, "Priority")
		},
	}

	templateExecutor := template.Must(template.New("").
		Funcs(funcMap).
		Parse(readmeTemplate))

	templateExecutor.Execute(f, struct {
		Rules []validations.Rule
	}{allRules})
}

var readmeTemplate = `# Rules
{{- range $rule := .Rules }}
*  [{{ $rule.LintID }}](#{{ lower $rule.LintID }})
{{- end }}

{{ range $rule := .Rules }}
## {{ $rule.LintID }}

> _{{ $rule.Summary }}_

{{ $rule.Details }}

Priority: **{{ shortPriority $rule.Priority.String }}**  
Analyzes: {{ range $command := $rule.Commands }}<kbd><a href="https://docs.docker.com/engine/reference/builder/#{{ $command }}">{{ $command.Upper }}</a></kbd> {{- end }}
{{ end }}
`
