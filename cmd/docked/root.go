package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/sirupsen/logrus"
)

// Build param: version
var version = "0.0.0"

// Build param: date
var date = "1970-01-01"

// Build param: commit
var commit = ""

// Build param: projectName
var projectName = "docked"

// CLI defines the command-line interface
var CLI struct {
	Config string `help:"Config file (default is $HOME/.docked.yaml)" type:"path"`

	Analyze AnalyzeCmd `cmd:"" help:"Analyze a Dockerfile for issues"`

	Version kong.VersionFlag `short:"v" help:"Print version information"`
}

func main() {
	initLogging()

	formattedVersion := fmt.Sprintf("%s (%s) %s", version, commit, date)

	ctx := kong.Parse(&CLI,
		kong.Name(projectName),
		kong.Description(`Dockerfile linting tool which aims to pull many
best practices and recommendations from multiple sources:

  * OWASP
  * Docker Official Documentation
  * Community recommendations
  * Package manager bug trackers`),
		kong.UsageOnError(),
		kong.Vars{
			"version": formattedVersion,
		},
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

// initLogging initializes logging used by the tool.
func initLogging() {
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info"
	}
	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		ll = logrus.ErrorLevel
	}
	logrus.SetLevel(ll)
	logrus.SetOutput(os.Stderr)
}
