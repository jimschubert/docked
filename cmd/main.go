package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jimschubert/docked"
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/validations"
	"github.com/urfave/cli/v2"

	"github.com/sirupsen/logrus"
)

var version = "0.0.0"
var date = "1970-01-01"
var commit = ""
var projectName = ""

func main() {
	initLogging()
	buildVersion := fmt.Sprintf("%s (%s) %s", version, commit, date)
	app := &cli.App{
		Usage:   "make an explosive entrance",
		Version: buildVersion,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dockerfile",
				Value: "./Dockerfile",
				Usage: "Path to dockerfile (defaults to ./Dockerfile)",
			},
			&cli.StringSliceFlag{
				Name:  "ignore",
				Aliases: []string{"i"},
				Usage: "The lint options to ignore",
			},
			&cli.BoolFlag{
				Name:  "no-buildkit-warnings",
				Usage: "Whether to suppress Docker parser warnings",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				TakesFile: true,
			},
		},
		Action: func(c *cli.Context) error {
			dockerfileOption := c.String("dockerfile")
			passedIgnores := c.StringSlice("ignore")
			customConfigPath := c.String("config")

			if len(dockerfileOption) == 0 {
				logrus.Fatal("No Dockerfile location provided")
			}

			config := buildConfig(passedIgnores, customConfigPath)

			application := docked.Docked{
				Config:                   config,
				SuppressBuildKitWarnings: c.Bool("no-buildkit-warnings"),
			}
			results, err := application.Analyze(dockerfileOption)
			if err != nil {
				return err
			}
			if len(results.Evaluated) == 0 {
				logrus.Warning("No validations selected")
			}
			printValidationResults(results.Evaluated)
			printRulesSkipped(results.NotEvaluated)

			errorCount := 0
			for _, validation := range results.Evaluated {
				if validation.ValidationResult.Result == model.Failure {
					errorCount += 1
				}
			}
			if errorCount > 0 {
				logrus.Fatalf("There were %d validation failures", errorCount)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.WithError(err).Fatalf("execution failed.")
	}
}

func buildConfig(passedIgnores []string, customConfigPath string) docked.Config {
	config := docked.Config{}
	if len(passedIgnores) > 0 {
		config.Ignore = passedIgnores
	}
	if len(customConfigPath) > 0 {
		err := config.Load(customConfigPath)
		if err != nil {
			logrus.WithError(err).Fatal("Unable to load the config path specified")
		}
	}
	return config
}

func printRulesSkipped(validations []validations.Validation) {
	for _, v := range validations {
		indicator := "#"
		priority := strings.TrimRight((*v.Rule).Priority().String(), "Priority")
		logrus.Printf("%s %-8s %s \n\t%s", indicator, priority, v.ID, v.Details)
	}
}

func printValidationResults(validations []validations.Validation) {
	for _, v := range validations {
		var indicator string
		priority := strings.TrimRight((*v.Rule).Priority().String(), "Priority")
		if v.ValidationResult.Result == model.Success {
			indicator = "✔"
			logrus.Printf("%s %-8s %s \n\t%s> %s\n\t%s", indicator, priority, v.ID, v.Contexts[0].Locations, v.Contexts[0].Line, v.Details)
		} else {
			indicator = "⨯"
			logrus.Errorf("%s %-8s %s \n\t%s> %s\n\t%s", indicator, priority, v.ID, v.Contexts[0].Locations, v.Contexts[0].Line, v.Details)
		}
	}
}

func initLogging() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info"
	}
	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)
	logrus.SetOutput(os.Stderr)
}
