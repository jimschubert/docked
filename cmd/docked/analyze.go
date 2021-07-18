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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type analyzeCommandOptions struct {
	Dockerfile         string
	NoBuildKitWarnings bool
	Ignores            []string
	ReportingType      string
	RegexEngine        string
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

func configureRegexEngine(opts analyzeCommandOptions) {
	switch opts.RegexEngine {
	case "regexp":
		model.SetRegexEngine(model.RegexpEngine)
	case "regexp2":
		fallthrough
	default:
		model.SetRegexEngine(model.Regexp2Engine)
	}
}

// newAnalyzeCommand creates the cobra.Command object with closed opts for use via `docked analyze`
func newAnalyzeCommand() *cobra.Command {
	opts := analyzeCommandOptions{
		Dockerfile:         "./Dockerfile",
		NoBuildKitWarnings: false,
		Ignores:            []string{},
		ReportingType:      "text",
		RegexEngine:        "regexp2",
	}

	analyzeCmd := &cobra.Command{
		Use:   "analyze [FILE]",
		Short: "Analyze a Dockerfile for issues",
		Long: `Analyze a Dockerfile for issues
If not provided, FILE defaults to ./Dockerfile
`,
		Args:       cobra.MaximumNArgs(1),
		ArgAliases: []string{"FILE"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.Dockerfile = args[0]
			}

			configureRegexEngine(opts)

			config := buildConfig(opts.Ignores, cfgFile)

			application := docked.Docked{
				Config:                   config,
				SuppressBuildKitWarnings: opts.NoBuildKitWarnings,
			}

			results, err := application.Analyze(opts.Dockerfile)
			cobra.CheckErr(err)

			if len(results.Evaluated) == 0 {
				logrus.Warning("No validations selected")
			}

			switch opts.ReportingType {
			case "json":
				var out bytes.Buffer
				b, err := json.Marshal(results)
				cobra.CheckErr(err)
				err = json.Indent(&out, b, "", "  ")
				cobra.CheckErr(err)
				_, _ = fmt.Fprintf(os.Stdout, "%s", out.Bytes())
			case "html":
				r := reporter.HTMLReporter{
					DockerfilePath: opts.Dockerfile,
				}
				err = r.Write(results)
				cobra.CheckErr(err)

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
				cobra.CheckErr(err)
			}

			errorCount := 0
			for _, validation := range results.Evaluated {
				if validation.ValidationResult.Result == model.Failure {
					errorCount++
				}
			}
			if errorCount > 0 {
				os.Exit(1)
			}
		},
	}

	analyzeCmd.Flags().BoolVarP(&opts.NoBuildKitWarnings, "no-buildkit-warnings", "k", opts.NoBuildKitWarnings, "Whether to suppress Docker parser warnings")
	analyzeCmd.Flags().StringSliceVarP(&opts.Ignores, "ignore", "i", opts.Ignores, "The lint ids to ignore")
	analyzeCmd.Flags().StringVarP(&opts.ReportingType, "report-type", "", opts.ReportingType, "The type of reporting output (text, json)")
	analyzeCmd.Flags().StringVarP(&opts.RegexEngine, "regex-engine", "", opts.RegexEngine, "The regex engine to use (regexp, regexp2)")
	viper.SetDefault("no-buildkit-warnings", opts.NoBuildKitWarnings)
	viper.SetDefault("ignore", opts.Ignores)
	viper.SetDefault("report-type", opts.ReportingType)
	viper.SetDefault("regex-engine", opts.RegexEngine)

	return analyzeCmd
}

func init() {
	analyzeCmd := newAnalyzeCommand()
	rootCmd.AddCommand(analyzeCmd)
}
