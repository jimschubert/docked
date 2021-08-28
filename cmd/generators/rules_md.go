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
				id := rule.GetLintID()
				if !seenRule[id] {
					allRules = append(allRules, rule)
					seenRule[id] = true
				}
			}
		}
	}

	sort.Slice(allRules, func(i, j int) bool {
		return allRules[i].GetLintID() < allRules[j].GetLintID()
	})

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"shortPriority": func(s string) string {
			return strings.TrimSuffix(s, "Priority")
		},
		"anchorText": func(s string) string {
			return strings.ReplaceAll(s, ":", "")
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
*  [{{ $rule.GetLintID }}](#{{ anchorText (lower $rule.GetLintID) }})
{{- end }}

{{ range $rule := .Rules }}
## {{ $rule.GetLintID }}

> _{{ $rule.GetSummary | html }}_

{{ $rule.GetDetails | html }}

Priority: **{{ shortPriority $rule.GetPriority.String }}**  
Analyzes: {{ range $command := $rule.GetCommands }}<kbd><a href="https://docs.docker.com/engine/reference/builder/#{{ $command }}">{{ $command.Upper }}</a></kbd> {{- end }}
{{ end }}
`
