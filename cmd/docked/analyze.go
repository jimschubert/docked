package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jimschubert/docked"
	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/reporter"
	"github.com/sirupsen/logrus"
)

// AnalyzeCmd represents the analyze command
type AnalyzeCmd struct {
	File               string   `arg:"" optional:"" type:"path" default:"./Dockerfile" help:"Dockerfile to analyze (default: ./Dockerfile)"`
	NoBuildKitWarnings bool     `short:"k" help:"Suppress Docker parser warnings"`
	Ignore             []string `short:"i" help:"Lint IDs to ignore"`
	ReportType         string   `enum:"text,json,html" default:"text" help:"Report output type (text, json, html)"`
	RegexEngine        string   `enum:"regexp,regexp2" default:"regexp2" help:"Regex engine to use (regexp, regexp2)"`
}

// Run executes the analyze command
func (a *AnalyzeCmd) Run() error {
	// Configure regex engine
	switch a.RegexEngine {
	case "regexp":
		model.SetRegexEngine(model.RegexpEngine)
	case "regexp2":
		fallthrough
	default:
		model.SetRegexEngine(model.Regexp2Engine)
	}

	// Build config
	config := docked.Config{}
	if len(a.Ignore) > 0 {
		config.Ignore = a.Ignore
	}
	if len(CLI.Config) > 0 {
		err := config.Load(CLI.Config)
		if err != nil {
			logrus.WithError(err).Fatal("Unable to load the config path specified")
		}
	}

	// Create application
	application := docked.Docked{
		Config:                   config,
		SuppressBuildKitWarnings: a.NoBuildKitWarnings,
	}

	// Analyze
	results, err := application.Analyze(a.File)
	if err != nil {
		return err
	}

	if len(results.Evaluated) == 0 {
		logrus.Warning("No validations selected")
	}

	// Generate report
	switch a.ReportType {
	case "json":
		var out bytes.Buffer
		b, err := json.Marshal(results)
		if err != nil {
			return err
		}
		err = json.Indent(&out, b, "", "  ")
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(os.Stdout, "%s", out.Bytes())
	case "html":
		r := reporter.HTMLReporter{
			DockerfilePath: a.File,
		}
		err = r.Write(results)
		if err != nil {
			return err
		}

		if absPath, err := filepath.Abs(r.OutDirectory); err == nil {
			fmt.Printf("HTML was output to: %s", absPath)
		}
	case "text":
		fallthrough
	default:
		r := reporter.TextReporter{
			DisableColors: false,
			Out:           os.Stdout,
		}
		err = r.Write(results)
		if err != nil {
			return err
		}
	}

	// Check for errors and exit with non-zero if found
	errorCount := 0
	for _, validation := range results.Evaluated {
		if validation.ValidationResult.Result == model.Failure {
			errorCount++
		}
	}
	if errorCount > 0 {
		os.Exit(1)
	}

	return nil
}
