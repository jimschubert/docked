package docked

import (
	"testing"

	"github.com/jimschubert/docked/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig_Deserialize(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name:    "contains ignore only",
			args:    args{helperTestData(t, "config/ignore_only.yml")},
			want:    Config{Ignore: []string{"D5:secret-aws-access-key"}},
			wantErr: false,
		},
		{
			name: "contains rule overrides only (as array of objects)",
			args: args{helperTestData(t, "config/rules_only.yml")},
			want: Config{RuleOverrides: &RuleOverrides{
				{"D5:secret-aws-access-key", model.LowPriority.Ptr()},
				{"D5:secret-aws-secret-access-key", model.CriticalPriority.Ptr()},
			}},
			wantErr: false,
		},
		{
			name: "contains rule overrides only (as map)",
			args: args{helperTestData(t, "config/rules_kvp.yml")},
			want: Config{RuleOverrides: &RuleOverrides{
				{"D5:secret-aws-access-key", model.LowPriority.Ptr()},
				{"D5:secret-aws-secret-access-key", model.CriticalPriority.Ptr()},
			}},
			wantErr: false,
		},
		{
			name: "contains ignores and rule overrides",
			args: args{helperTestData(t, "config/ignores_and_rules.yml")},
			want: Config{
				Ignore: []string{"D5:secret-aws-access-key", "D5:secret-aws-secret-access-key"},
				RuleOverrides: &RuleOverrides{
					{"D7:tagged-latest", model.CriticalPriority.Ptr()},
					{"D7:tagged-latest-builder", model.HighPriority.Ptr()},
					{"DC:consider-multistage", model.CriticalPriority.Ptr()},
				}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{}
			err := yaml.Unmarshal(tt.args.b, &c)
			if (err != nil) != tt.wantErr {
				t.Errorf("yaml.Unmarshal to Config{} error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// TODO: Sort these configs to avoid random test failures on the map test (seems to come from yaml.Node)
			assert.Equal(t, tt.want, c, "yaml.Unmarshal to  got = %v, want %v", c, tt.want)
		})
	}
}
