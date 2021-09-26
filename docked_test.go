package docked

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/rules"
	"github.com/jimschubert/docked/model/validations"
)

func printEvaluated(evaluated []validations.Validation) {
	for _, r := range evaluated {
		b := bytes.Buffer{}
		_, _ = b.WriteString(fmt.Sprintf("%s - %s", r.ID, r.Result.String()))
		if r.Result == model.Failure {
			for _, c := range r.Contexts {
				if c.CausedFailure {
					_, _ = b.WriteString(fmt.Sprintf(" * [%2d] %s\n", c.Locations[0].Start.Line, c.Line))
				}
			}
		} else {
			_, _ = b.WriteString("\n")
		}
		_, _ = fmt.Fprint(os.Stdout, b.String())
	}
}

// ExampleDocked_Analyze provides an example of programmatically invoking Docked.Analyze with default rules
func ExampleDocked_Analyze() {
	c := Config{}
	if err := c.Load("./testdata/config/example.yml"); err != nil {
		panic(err)
	}

	d := Docked{
		Config:                   c,
		SuppressBuildKitWarnings: true,
	}

	result, err := d.Analyze("./testdata/minimal.dockerfile")
	if err != nil {
		panic("Failed to analyze dockerfile")
	}

	// programmatically consume array of evaluated and/or not-evaluated rules
	printEvaluated(result.Evaluated)

	// Output:
	// D5:no-debian-frontend - Success
	// D5:secret-aws-access-key - Success
	// D5:secret-aws-secret-access-key - Success
	// DC:avoid-sudo - Success
	// DC:consider-multistage - Success
	// DC:curl-without-fail - Success
	// DC:gpg-without-batch - Success
	// DC:gpg-without-batch - Success
	// DC:layered-ownership-change - Success
}

// ExampleDocked_Analyze_withCustomRules provides an example of programmatically invoking Docked.Analyze with custom rules
func ExampleDocked_Analyze_withCustomRules() {
	c := Config{}
	// The config file will define a rule named adding-full-directory
	if err := c.Load("./testdata/config/example_custom.yml"); err != nil {
		panic(err)
	}

	d := Docked{
		Config:                   c,
		SuppressBuildKitWarnings: true,
	}

	result, err := d.Analyze("./testdata/minimal_custom.dockerfile")
	if err != nil {
		panic("Failed to analyze dockerfile")
	}

	// programmatically consume array of evaluated and/or not-evaluated rules
	printEvaluated(result.Evaluated)

	// Output:
	// D0:adding-full-directory - Failure * [ 7] ADD . /go/src/app
	// D5:no-debian-frontend - Success
	// D5:secret-aws-access-key - Success
	// D5:secret-aws-secret-access-key - Success
	// DC:avoid-sudo - Success
	// DC:consider-multistage - Success
	// DC:curl-without-fail - Success
	// DC:gpg-without-batch - Success
	// DC:gpg-without-batch - Success
	// DC:layered-ownership-change - Success
}

// ExampleDocked_AnalyzeWithRuleList provides an example of programmatically invoking Docked.AnalyzeWithRuleList
// with user-defined rules. See also reporter.TextReporter and reporter.HTMLReporter for in-built output formatters.
func ExampleDocked_AnalyzeWithRuleList() {
	d := Docked{}

	// user can extend default rule set or define their own
	activeRules := rules.RuleList{}
	myRule := validations.SimpleRegexRule{
		Name:     "no-distroless",
		Pattern:  `\bgcr\.io/distroless\b`,
		Priority: model.CriticalPriority,
		Command:  commands.From,
	}
	activeRules.AddRule(myRule)

	result, err := d.AnalyzeWithRuleList("./testdata/minimal.dockerfile", ConfiguredRules{Active: activeRules})
	if err != nil {
		panic("Failed to analyze dockerfile")
	}

	// programmatically consume array of evaluated and/or not-evaluated rules
	printEvaluated(result.Evaluated)

	// Output:
	// D7:no-distroless - Success
	// D7:no-distroless - Failure * [13] FROM gcr.io/distroless/base-debian10
}

func v(name string, result model.Valid) validations.Validation {
	return validations.Validation{ID: name, ValidationResult: validations.ValidationResult{Result: result}}
}

func singleValidationSlice(name string, result model.Valid) []validations.Validation {
	return []validations.Validation{
		v(name, result),
	}
}

func TestDocked_buildConfiguredRules_ignored(t *testing.T) {
	type args struct {
		config   Config
		location string
	}
	tests := []struct {
		name            string
		args            args
		ignoredContains []validations.Validation
		wantErr         bool
	}{
		{
			name: "ignoring rules",
			args: args{
				config:   Config{Ignore: []string{"D2:single-cmd", "D3:avoid-copy-all", "DC:curl-without-fail"}},
				location: "./testdata/minimal.dockerfile",
			},
			ignoredContains: []validations.Validation{
				v("D2:single-cmd", model.Ignored),
				v("D3:avoid-copy-all", model.Ignored),
				v("DC:curl-without-fail", model.Ignored),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.ignoredContains) == 0 {
				// need to force this check since this method will validate in iteration
				t.Errorf("This test is written incorrectly and failed pre-check")
				return
			}

			expected := make(map[string]validations.Validation)
			actual := make(map[string]bool)
			configuredRules := buildConfiguredRules(tt.args.config)

			for _, contain := range tt.ignoredContains {
				expected[contain.ID] = contain
			}

			for _, r := range configuredRules.Inactive {
				for _, rule := range *r {
					actual[rule.GetLintID()] = true
				}
			}
			for _, ignored := range expected {
				if !actual[ignored.ID] {
					t.Errorf("Expected '%v' to be ignored", ignored.ID)
				}
			}
		})
	}
}

func TestDocked_AnalyzeWithRuleList(t *testing.T) {
	type args struct {
		config   Config
		location string
	}
	tests := []struct {
		name    string
		args    args
		want    AnalysisResult
		wantErr bool
	}{
		// region avoid-copy-all
		{
			name: "avoid-copy-all [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D3:avoid-copy-all"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D3:avoid-copy-all", model.Success)},
		},
		{
			name: "avoid-copy-all (recommendation)",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D3:avoid-copy-all"}},
				location: "./testdata/copy_all.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D3:avoid-copy-all", model.Recommendation)},
		},
		{
			name: "avoid-copy-all (not in indexed builder context)",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D3:avoid-copy-all"}},
				location: "./testdata/copy_all_indexed_builder.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D3:avoid-copy-all", model.Success)},
		},
		{
			name: "avoid-copy-all (not in named builder context)",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D3:avoid-copy-all"}},
				location: "./testdata/copy_all_named_builder.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D3:avoid-copy-all", model.Success)},
		},
		// endregion avoid-copy-all

		// region avoid-sudo
		{
			name: "avoid-sudo",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:avoid-sudo"}},
				location: "./testdata/avoid_sudo.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:avoid-sudo", model.Recommendation)},
		},
		{
			name: "avoid-sudo [su] (recommendation)",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:avoid-sudo"}},
				location: "./testdata/avoid_sudo_su.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:avoid-sudo", model.Recommendation)},
		},
		{
			name: "avoid-sudo [gosu]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:avoid-sudo"}},
				location: "./testdata/avoid_sudo_gosu.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:avoid-sudo", model.Success)},
		},
		{
			name: "avoid-sudo [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:avoid-sudo"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:avoid-sudo", model.Success)},
		},
		// endregion avoid-sudo

		// region consider-multistage
		{
			name: "consider-multistage [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:consider-multistage"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:consider-multistage", model.Success)},
		},
		{
			name: "consider-multistage [mvn]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:consider-multistage"}},
				location: "./testdata/consider_multistage.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:consider-multistage", model.Failure)},
		},
		{
			name: "consider-multistage [go]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:consider-multistage"}},
				location: "./testdata/consider_multistage_go_build.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:consider-multistage", model.Failure)},
		},
		// endregion consider-multistage

		// region curl-without-fail
		{
			name: "curl-without-fail",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:curl-without-fail"}},
				location: "./testdata/curl_without_fail.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:curl-without-fail", model.Failure)},
		},
		{
			name: "curl-without-fail [issue #2]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:curl-without-fail"}},
				location: "./testdata/curl_without_fail_issue2.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:curl-without-fail", model.Success)},
		},
		{
			name: "curl-without-fail [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:curl-without-fail"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:curl-without-fail", model.Success)},
		},
		// endregion curl-without-fail

		// region gpg-without-batch
		{
			name: "gpg-without-batch",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:gpg-without-batch"}},
				location: "./testdata/gpg_without_batch.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:gpg-without-batch", model.Failure)},
		},
		{
			name: "gpg-without-batch [batch]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:gpg-without-batch"}},
				location: "./testdata/gpg_with_batch.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:gpg-without-batch", model.Success)},
		},
		{
			name: "gpg-without-batch [no-tty]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:gpg-without-batch"}},
				location: "./testdata/gpg_with_no_tty.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:gpg-without-batch", model.Success)},
		},
		{
			name: "gpg-without-batch [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:gpg-without-batch"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: []validations.Validation{
				v("DC:gpg-without-batch", model.Success),
				v("DC:gpg-without-batch", model.Success),
			}},
		},
		// endregion gpg-without-batch

		// region maintainer-deprecated
		{
			name: "maintainer-deprecated",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DA:maintainer-deprecated"}},
				location: "./testdata/maintainer.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DA:maintainer-deprecated", model.Failure)},
		},
		{
			name: "maintainer-deprecated [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DA:maintainer-deprecated"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{NotEvaluated: singleValidationSlice("DA:maintainer-deprecated", model.Skipped)},
		},
		// endregion maintainer-deprecated

		// region no-debian-frontend
		{
			name: "no-debian-frontend [env]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D5:no-debian-frontend"}},
				location: "./testdata/debian_frontend_env.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D5:no-debian-frontend", model.Failure)},
		},
		{
			name: "no-debian-frontend [arg]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D5:no-debian-frontend"}},
				location: "./testdata/debian_frontend_arg.dockerfile",
			},
			want: AnalysisResult{NotEvaluated: singleValidationSlice("D5:no-debian-frontend", model.Success)},
		},
		{
			name: "no-debian-frontend [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D5:no-debian-frontend"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D5:no-debian-frontend", model.Success)},
		},
		// endregion no-debian-frontend

		// region tagged-latest
		{
			name: "tagged-latest [scratch]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D7:tagged-latest"}},
				location: "./testdata/scratch.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D7:tagged-latest", model.Success)},
		},
		{
			name: "tagged-latest [digest]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D7:tagged-latest"}},
				location: "./testdata/image_with_digest.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D7:tagged-latest", model.Success)},
		},
		{
			name: "tagged-latest [final image]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D7:tagged-latest"}},
				location: "./testdata/tagged_latest.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D7:tagged-latest", model.Failure)},
		},
		{
			name: "tagged-latest [builder image]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D7:tagged-latest"}},
				location: "./testdata/tagged_latest_builder.dockerfile",
			},
			want: AnalysisResult{
				Evaluated:    singleValidationSlice("D7:tagged-latest", model.Success),
				NotEvaluated: singleValidationSlice("D7:tagged-latest", model.Skipped),
			},
		},
		// endregion tagged-latest

		// region tagged-latest-builder
		{
			name: "tagged-latest-builder [scratch]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D7:tagged-latest-builder"}},
				location: "./testdata/scratch.dockerfile",
			},
			want: AnalysisResult{NotEvaluated: singleValidationSlice("D7:tagged-latest-builder", model.Success)},
		},
		{
			name: "tagged-latest-builder [final image]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D7:tagged-latest-builder"}},
				location: "./testdata/tagged_latest.dockerfile",
			},
			want: AnalysisResult{NotEvaluated: singleValidationSlice("D7:tagged-latest-builder", model.Success)},
		},
		{
			name: "tagged-latest-builder [builder image]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D7:tagged-latest-builder"}},
				location: "./testdata/tagged_latest_builder.dockerfile",
			},
			want: AnalysisResult{
				Evaluated:    singleValidationSlice("D7:tagged-latest-builder", model.Failure),
				NotEvaluated: singleValidationSlice("D7:tagged-latest-builder", model.Skipped),
			},
		},
		// endregion tagged-latest

		// region secret-aws-access-key
		{
			name: "secret-aws-access-key [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D5:secret-aws-access-key"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D5:secret-aws-access-key", model.Success)},
		},
		{
			name: "secret-aws-access-key",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D5:secret-aws-access-key"}},
				location: "./testdata/aws_secrets.dockerfile",
			},
			want: AnalysisResult{Evaluated: []validations.Validation{
				v("D5:secret-aws-access-key", model.Failure),
				v("D5:secret-aws-access-key", model.Success),
			}},
		},
		// endregion secret-aws-access-key

		// region secret-aws-access-key
		{
			name: "secret-aws-secret-access-key [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D5:secret-aws-secret-access-key"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D5:secret-aws-secret-access-key", model.Success)},
		},
		{
			name: "secret-aws-secret-access-key",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D5:secret-aws-secret-access-key"}},
				location: "./testdata/aws_secrets.dockerfile",
			},
			want: AnalysisResult{Evaluated: []validations.Validation{
				v("D5:secret-aws-secret-access-key", model.Success),
				v("D5:secret-aws-secret-access-key", model.Failure),
			}},
		},
		// endregion secret-aws-access-key

		// region questionable-expose
		{
			name: "questionable-expose",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D6:questionable-expose"}},
				location: "./testdata/questionable_expose.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D6:questionable-expose", model.Failure)},
		},
		{
			name: "questionable-expose [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D6:questionable-expose"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{NotEvaluated: singleValidationSlice("D6:questionable-expose", model.Skipped)},
		},
		// endregion questionable-expose

		// region layered-ownership-change
		{
			name: "layered-ownership-change",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:layered-ownership-change"}},
				location: "./testdata/layered_ownership_change.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:layered-ownership-change", model.Failure)},
		},
		{
			name: "questionable-expose [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:layered-ownership-change"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:layered-ownership-change", model.Success)},
		},
		// endregion questionable-expose

		// region single-cmd
		{
			name: "single-cmd",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D2:single-cmd"}},
				location: "./testdata/single_cmd.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D2:single-cmd", model.Failure)},
		},
		{
			name: "single-cmd [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D2:single-cmd"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{NotEvaluated: singleValidationSlice("D2:single-cmd", model.Success)},
		},
		// endregion questionable-expose

		// region oci-labels
		{
			name: "oci-labels",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D9:oci-labels"}},
				location: "./testdata/oci_labels.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("D9:oci-labels", model.Success)},
		},
		{
			name: "oci-labels [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"D9:oci-labels"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{NotEvaluated: singleValidationSlice("D9:oci-labels", model.Recommendation)},
		},
		// endregion oci-labels

		// region minimize-layers
		{
			name: "minimize-layers",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:minimize-layers"}},
				location: "./testdata/minimize_layers.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:minimize-layers", model.Recommendation)},
		},
		{
			name: "minimize-layers [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:minimize-layers"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:minimize-layers", model.Success)},
		},
		// endregion minimize-layers

		// region sort-installer-args
		{
			name: "sort-installer-args [apt (sorted)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_sorted.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Success)},
		},
		{
			name: "sort-installer-args [apt (unsorted)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_unsorted.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Recommendation)},
		},
		{
			name: "sort-installer-args [apt-get (sorted)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_get_sorted.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Success)},
		},
		{
			name: "sort-installer-args [apt-get (unsorted)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_get_unsorted.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Recommendation)},
		},
		{
			name: "sort-installer-args [apt-get (switchless)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_get_switchless.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Recommendation)},
		},
		{
			name: "sort-installer-args [apt-get (apt_get_with_command.dockerfile)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_get_with_command.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Recommendation)},
		},
		{
			name: "sort-installer-args [npm (sorted)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/npm_sorted.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Success)},
		},
		{
			name: "sort-installer-args [npm (unsorted)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/npm_unsorted.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Recommendation)},
		},
		{
			name: "sort-installer-args [npm (switchless)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/npm_switchless.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Success)},
		},
		{
			name: "sort-installer-args [yum (sorted)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_get_sorted.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Success)},
		},
		{
			name: "sort-installer-args [yum (unsorted)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_get_unsorted.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Recommendation)},
		},
		{
			name: "sort-installer-args [yum (switchless)]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/sort_installer_args/apt_get_switchless.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Recommendation)},
		},
		{
			name: "sort-installer-args [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:sort-installer-args"}},
				location: "./testdata/minimal.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:sort-installer-args", model.Success)},
		},
		// endregion minimize-layers
		// region apt-get-update-install
		{
			name: "apt-get-update-install [separate]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:apt-get-update-install"}},
				location: "./testdata/apt_get_best_practices/install_separate_instruction.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:apt-get-update-install", model.Failure)},
		},
		{
			name: "apt-get-update-install [same]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:apt-get-update-install"}},
				location: "./testdata/apt_get_best_practices/install_same_instruction.dockerfile",
			},
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:apt-get-update-install", model.Success)},
		},
		{
			name: "apt-get-update-install [minimal]",
			args: args{
				config:   Config{SkipDefaultRules: true, IncludeRules: []string{"DC:apt-get-update-install"}},
				location: "./testdata/minimal.dockerfile",
			},
			// minimal.dockerfile is alpine, so apt-get specific reports as "skipped"
			want: AnalysisResult{Evaluated: singleValidationSlice("DC:apt-get-update-install", model.Skipped)},
		},
		// endregion minimize-layers
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Docked{Config: tt.args.config, SuppressBuildKitWarnings: true}
			configuredRules := buildConfiguredRules(tt.args.config)
			got, err := d.AnalyzeWithRuleList(tt.args.location, configuredRules)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeWithRuleList error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			showError := func() {
				t.Errorf("AnalyzeWithRuleList() = %#v, want (w/o contexts) %#v", got, tt.want)
			}

			// We don't use  reflect.DeepEqual here to simplify building out tests (no path, full rule definition, etc)
			if len(got.Evaluated) != len(tt.want.Evaluated) {
				showError()
				return
			}

			if len(got.NotEvaluated) != len(tt.want.NotEvaluated) {
				showError()
				return
			}

			for i, actualValidation := range got.Evaluated {
				expectedValidation := tt.want.Evaluated[i]

				if actualValidation.ID != expectedValidation.ID {
					t.Logf("Actual ID = %s, expected ID = %s", actualValidation.ID, expectedValidation.ID)
					showError()
					return
				}

				if actualValidation.Result != expectedValidation.Result {
					t.Logf("Actual result = %#v, expected result = %#v", actualValidation.Result, expectedValidation.Result)
					showError()
					return
				}

				if len(expectedValidation.Details) == 0 {
					t.Log("expectedValidation.Details not specified, not checking…")
				} else if actualValidation.Details != expectedValidation.Details {
					t.Logf("Actual Details = %s, expected Details = %s", actualValidation.Details, expectedValidation.Details)
					showError()
					return
				}

				if len(expectedValidation.Contexts) == 0 {
					t.Log("expectedValidation.Contexts not specified, not checking…")
				} else if reflect.DeepEqual(actualValidation.Contexts, expectedValidation.Contexts) {
					t.Logf("Actual Contexts = %v, expected Contexts = %v", actualValidation.Details, expectedValidation.Details)
					showError()
					return
				}
			}
		})
	}
}
