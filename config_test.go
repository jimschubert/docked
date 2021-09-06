package docked

import (
	"testing"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/validations"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name:    "contains ignore only",
			args:    args{"testdata/config/ignore_only.yml"},
			want:    Config{Ignore: []string{"D5:secret-aws-access-key"}},
			wantErr: false,
		},
		{
			name: "contains rule overrides only (as array of objects)",
			args: args{"testdata/config/rules_only.yml"},
			want: Config{RuleOverrides: &RuleOverrides{
				{"D5:secret-aws-access-key", model.LowPriority.Ptr()},
				{"D5:secret-aws-secret-access-key", model.CriticalPriority.Ptr()},
			}},
			wantErr: false,
		},
		{
			name: "contains rule overrides only (as map)",
			args: args{"testdata/config/rules_kvp.yml"},
			want: Config{RuleOverrides: &RuleOverrides{
				{"D5:secret-aws-access-key", model.LowPriority.Ptr()},
				{"D5:secret-aws-secret-access-key", model.CriticalPriority.Ptr()},
			}},
			wantErr: false,
		},
		{
			name: "contains ignores and rule overrides",
			args: args{"testdata/config/ignores_and_rules.yml"},
			want: Config{
				Ignore: []string{"D5:secret-aws-access-key", "D5:secret-aws-secret-access-key"},
				RuleOverrides: &RuleOverrides{
					{"D7:tagged-latest", model.CriticalPriority.Ptr()},
					{"D7:tagged-latest-builder", model.HighPriority.Ptr()},
					{"DC:consider-multistage", model.CriticalPriority.Ptr()},
				}},
			wantErr: false,
		},
		{
			name: "contains ignores and rule overrides",
			args: args{"testdata/config/full_with_custom_rules.yml"},
			want: Config{
				Ignore: []string{"D5:secret-aws-access-key", "D5:secret-aws-secret-access-key"},
				RuleOverrides: &RuleOverrides{
					{"D7:tagged-latest", model.CriticalPriority.Ptr()},
					{"D7:tagged-latest-builder", model.HighPriority.Ptr()},
					{"DC:consider-multistage", model.CriticalPriority.Ptr()},
				},
				CustomRules: []validations.SimpleRegexRule{
					{
						Name:     "no funny business",
						Summary:  "Prevent common typo on our team",
						Details:  "Jim keeps mistyping rm -rf",
						Pattern:  `rm -rf /\b`,
						Priority: model.CriticalPriority,
						Command:  commands.Run,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{}
			err := c.Load(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("yaml.Unmarshal to Config{} error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, c, "yaml.Unmarshal to  got = %v, want %v", c, tt.want)
		})
	}
}
