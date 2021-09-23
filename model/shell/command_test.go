package shell

import (
	"reflect"
	"testing"
)

func TestNewPosixCommand(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    []PosixCommand
		wantErr bool
	}{
		{
			name: "issue #2 - false positive for curl - should not produce curl command",
			args: args{input: "apk --no-cache update && apk --no-cache add curl"},
			want: []PosixCommand{
				{
					Name: "apk",
					Args: []string{"--no-cache", "update"},
				},
				{
					Name: "apk",
					Args: []string{"--no-cache", "add", "curl"},
				},
			},
		},
		{
			name: "issue #2 - false positive for curl - should produce curl command",
			args: args{input: "true && curl --fail https://example.com/file.json && apk --no-cache add curl"},
			want: []PosixCommand{
				{
					Name: "true",
					Args: []string{},
				},
				{
					Name: "curl",
					Args: []string{"--fail", "https://example.com/file.json"},
				},
				{
					Name: "apk",
					Args: []string{"--no-cache", "add", "curl"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPosixCommand(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPosixCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPosixCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}
