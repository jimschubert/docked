package rules

import (
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/shell"
	"github.com/jimschubert/docked/model/validations"
	log "github.com/sirupsen/logrus"
)

func aptGetUpdateInstall() validations.Rule {
	r := validations.MultiContextRule{
		Name:     "apt-get-update-install",
		Summary:  "You must perform apt-get update and install in same RUN layer",
		Details:  "Having apt-get update and install in separate RUN layers will break caching. Having install without update is not recommended. Include both commands in the same layer.",
		Priority: model.CriticalPriority,
		Commands: []commands.DockerCommand{commands.Run},
		URL:      model.StringPtr("https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#apt-get"),
		Evaluator: validations.MultiContextFullEvaluator{
			Fn: func(mcr *validations.MultiContextRule) *validations.ValidationResult {
				var updateIdx *int
				var installIdx *int
				result := model.Success
				for idx, nodeContext := range *mcr.ContextCache {
					posixCommands, err := shell.NewPosixCommandFromNode(&nodeContext.Node)
					if err != nil {
						log.Warnf("Unable to parse RUN command, validation not evaluated: %#v", nodeContext.Node.Location())
						return &validations.ValidationResult{
							Result:  model.Skipped,
							Details: mcr.GetSummary(),
						}
					}

					for _, command := range posixCommands {
						commandName := strings.TrimLeft(command.Name, `\`)
						if commandName == "apt-get" {
							current := idx
							switch command.Args[0] {
							case "install":
								installIdx = &current
							case "update":
								updateIdx = &current
							}
						}
					}
				}
				if installIdx == nil && updateIdx == nil {
					result = model.Skipped
				} else if installIdx != nil {
					if updateIdx == nil || (*installIdx > *updateIdx) {
						result = model.Failure
					}
				}
				return &validations.ValidationResult{
					Result:   result,
					Details:  mcr.GetSummary(),
					Contexts: *mcr.GetContexts(),
				}
			},
		},
	}
	return &r
}

func init() {
	AddRule(aptGetUpdateInstall())
}
