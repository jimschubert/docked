package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/shell"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

type gpgWithoutBatch struct {
}

func (g *gpgWithoutBatch) GetName() string {
	return "gpg-without-batch"
}

func (g *gpgWithoutBatch) GetSummary() string {
	return "GPG call without --batch (or --no-tty) may error."
}

func (g *gpgWithoutBatch) GetDetails() string {
	return "Running GPG without --batch (or --no-tty) may cause GPG to fail opening /dev/tty, resulting in docker build failures."
}

func (g *gpgWithoutBatch) GetPriority() model.Priority {
	return model.MediumPriority
}

func (g *gpgWithoutBatch) GetCommands() []commands.DockerCommand {
	return []commands.DockerCommand{commands.Run}
}

func (g *gpgWithoutBatch) GetCategory() *string {
	return nil
}

func (g *gpgWithoutBatch) GetURL() *string {
	return model.StringPtr("https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=913614")
}

func (g *gpgWithoutBatch) GetLintID() string {
	return validations.LintID(g)
}

func (g *gpgWithoutBatch) Evaluate(node *parser.Node, validationContext validations.ValidationContext) *validations.ValidationResult {
	trimStart := len(node.Value) + 1 // command plus trailing space
	commandText := node.Original[trimStart:]
	result := model.Success
	posixCommands, err := shell.NewPosixCommand(commandText)
	if err != nil {
		log.Warnf("Unable to parse RUN command, skipping validation for: %s", commandText)
		result = model.Skipped
	} else {
		for _, command := range posixCommands {
			if (command.Name == "gpg" || command.Name == `\gpg`) && command.Args != nil {
				var hasBatch bool
				var hasNoTty bool
				for _, arg := range command.Args {
					hasBatch = hasBatch || strings.Contains(arg, "--batch")
					hasNoTty = hasNoTty || strings.Contains(arg, "--no-tty")
				}

				if !hasBatch && !hasNoTty {
					result = model.Failure
					validationContext.CausedFailure = true
				}
			}
		}
	}

	return &validations.ValidationResult{
		Result:   result,
		Details:  g.GetSummary(),
		Contexts: []validations.ValidationContext{validationContext},
	}
}

func init() {
	r := gpgWithoutBatch{}
	AddRule(&r)
}
