package docked

import (
	"os"
	"path/filepath"

	"github.com/jimschubert/docked/model/docker"
	"github.com/jimschubert/docked/model/docker/command"
	"github.com/jimschubert/docked/model/rules"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

type Docked struct {
}

func (a *Docked) Analyze(location string) ([]validations.Validation, error) {
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

	rulesList := rules.DefaultRules()

	validationsRan := make([]validations.Validation, 0)

	//goland:noinspection ALL
	for _, node := range p.AST.Children {
		thisCommand := command.DockerCommand(node.Value)
		if rules, ok := rulesList[thisCommand]; ok {
			for _, rule := range *rules {
				result := rule.Evaluate(node)
				validationsRan = append(validationsRan, validations.Validation{
					ID:               rule.LintID(),
					ValidationResult: *result,
					Line:             node.Original,
					Range:            docker.FromParserRanges(node.Location()),
				})
			}
		}
	}

	return validationsRan, nil
}
