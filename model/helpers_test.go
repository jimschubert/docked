package model

import (
	"testing"
)

func TestStringSliceContains(t *testing.T) {
	type args struct {
		slice *[]string
		str   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "negative: zero-empty", args: args{&[]string{}, ""}, want: false},
		{name: "negative: zero-nonempty", args: args{&[]string{}, "asdf"}, want: false},
		{name: "negative: single-nonempty", args: args{&[]string{"plugh"}, "asdf"}, want: false},
		{name: "positive: single-nonempty", args: args{&[]string{"plugh"}, "plugh"}, want: true},
		{name: "positive: empty-empty", args: args{&[]string{""}, ""}, want: true},
		{name: "negative: nil slice", args: args{nil, ""}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContains(tt.args.slice, tt.args.str); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
