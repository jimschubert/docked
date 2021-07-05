package rules

import (
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
)

const (
	genericNoSecretsSummary = "Secrets should not be stored directly in the Dockerfile. You should remove and rotate any secrets used here."
)

func noAwsAccessKey() validations.Rule {
	// see https://aws.amazon.com/blogs/security/a-safer-way-to-distribute-aws-credentials-to-ec2/
	r := validations.SimpleRegexRule{
		Name:     "secret-aws-access-key",
		Summary:  genericNoSecretsSummary,
		Pattern:  `\bAK[A-Z0-9]{18}\b`,
		Priority: model.CriticalPriority,
		Command:  commands.Env,
	}
	return &r
}

func noAwsSecretAccessKey() validations.Rule {
	// see https://aws.amazon.com/blogs/security/a-safer-way-to-distribute-aws-credentials-to-ec2/
	r := validations.SimpleRegexRule{
		Name:     "secret-aws-secret-access-key",
		Summary:  genericNoSecretsSummary,
		Pattern:  `\b[A-Za-z0-9/+=]{40}\b`,
		Priority: model.CriticalPriority,
		Command:  commands.Env,
	}
	return &r
}

func init() {
	AddRule(noAwsAccessKey())
	AddRule(noAwsSecretAccessKey())
}
