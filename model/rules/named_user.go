package rules

import (
	"strconv"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"
)

func namedUser() validations.Rule {
	isUserName := func(name string) bool {
		i, err := strconv.Atoi(name)
		if _, ok := err.(*strconv.NumError); ok {
			return true
		}
		log.Debugf("user name string '%s' parsed as an integer: %d", name, i)
		return false
	}

	isUserInstructionValid := func(c *validations.NodeValidationContext) bool {
		i, err := instructions.ParseInstruction(&c.Node)
		if err != nil {
			log.Warnf("failed parsing user command: %s", err)
			return false
		}
		if user, ok := i.(*instructions.UserCommand); ok && user != nil {
			return isUserName(user.User)
		}

		return true
	}

	isCopyInstructionValid := func(c *validations.NodeValidationContext) bool {
		i, err := instructions.ParseInstruction(&c.Node)
		if err != nil {
			log.Warnf("failed parsing copy command: %s", err)
			return false
		}

		if copyCommand, ok := i.(*instructions.CopyCommand); ok && copyCommand != nil {
			chown := copyCommand.Chown
			if chown != "" {
				return isUserName(chown)
			}
		}

		return true
	}
	r := validations.MultiContextRule{
		Name:     "named-user",
		Summary:  "Reference a user by name rather than UID.",
		Details:  "Reference a user by name to avoid maintenance or runtime issues with generated IDs.",
		Priority: model.HighPriority,
		Commands: []commands.DockerCommand{commands.User, commands.Copy},
		URL:      model.StringPtr("https://devopsbootcamp.org/dockerfile-security-best-practices/#1-2-don-t-bind-to-a-specific-uid"),
		Evaluator: validations.MultiContextFullEvaluator{
			Fn: func(mcr *validations.MultiContextRule) *validations.ValidationResult {
				result := model.Success
				validationContexts := make([]validations.ValidationContext, 0)
				for _, nodeContext := range *mcr.ContextCache {
					thisCommand := commands.DockerCommand(nodeContext.Node.Value)
					var isValid bool
					switch thisCommand {
					case commands.User:
						isValid = isUserInstructionValid(&nodeContext)
					case commands.Copy:
						isValid = isCopyInstructionValid(&nodeContext)
					default:
						panic("named-user: Unexpected docker command during evaluation")
					}

					if !isValid {
						nodeContext.Context.CausedFailure = true
						result = model.Failure
					}
					validationContexts = append(validationContexts, nodeContext.Context)
				}
				return &validations.ValidationResult{
					Result:   result,
					Details:  mcr.GetSummary(),
					Contexts: validationContexts,
				}
			},
		},
	}
	return &r
}

func init() {
	n := namedUser()
	AddRule(n)
}
