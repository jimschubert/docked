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
			if len(dockerfileOption) == 0 {
				logrus.Fatal("No Dockerfile location provided")
			}
			config := docked.Config{}
			passedIgnores := c.StringSlice("ignore")
			if len(passedIgnores) > 0 {
				config.Ignore = passedIgnores
			}
			customConfigPath := c.String("config")
			if len(customConfigPath) > 0 {
				err := config.Load(customConfigPath)
				if err != nil {
					logrus.WithError(err).Fatal("Unable to load the config path specified")
				}
			}

			application := docked.Docked{
				Config:                   config,
				SuppressBuildKitWarnings: c.Bool("no-buildkit-warnings"),
			}
			validations, err := application.Analyze(dockerfileOption)
			if err != nil {
				return err
			}
			if len(validations) == 0 {
				logrus.Warning("No validations selected")
			}
			printValidationResults(validations)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.WithError(err).Fatalf("execution failed.")
	}
}

func printValidationResults(validations []validations.Validation) {
	var errCount = 0
	for _, v := range validations {
		if v.ValidationResult.Result == model.Failure {
			errCount += 1
		}
		var indicator string
		if v.ValidationResult.Result == model.Success {
			indicator = "✅"
		} else {
			indicator = "❌"
		}
		priority := strings.TrimRight((*v.Rule).Priority().String(), "Priority")
		logrus.Printf("%s %-8s: %s \n\t%s> %s\n\t%s", indicator, priority, v.ID, v.Contexts[0].Locations, v.Contexts[0].Line, v.Details)
	}
	if errCount > 0 {
		logrus.Fatalf("There were %d errors", errCount)
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
