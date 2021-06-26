package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

const (
	genericNoSecretsSummary = "Secrets should not be stored directly in the Dockerfile"
)

func noAwsAccessKey() validations.Rule {
	// see https://aws.amazon.com/blogs/security/a-safer-way-to-distribute-aws-credentials-to-ec2/
	return validations.NewSimpleRegexRule(
		"secret-aws-access-key",
		genericNoSecretsSummary,
		`\bAK[A-Z0-9]{18}\b`,
		model.HighPriority,
		commands.Env,
		nil,
		nil,
	)
}

func noAwsSecretAccessKey() validations.Rule {
	// see https://aws.amazon.com/blogs/security/a-safer-way-to-distribute-aws-credentials-to-ec2/
	return validations.NewSimpleRegexRule(
		"secret-aws-secret-access-key",
		genericNoSecretsSummary,
		`\b[A-Za-z0-9/+=]{40}\b`,
		model.HighPriority,
		commands.Env,
		nil,
		nil,
	)
}

func init() {
	AddRule(noAwsAccessKey())
	AddRule(noAwsSecretAccessKey())
}
