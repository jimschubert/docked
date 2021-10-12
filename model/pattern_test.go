package model

import (
	"reflect"
	"testing"
)

func TestNewPattern(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want Pattern
	}{
		{name: "simple", args: args{value: `.*`}, want: StandardPattern{pattern: `.*`, engine: goRegexpEngine{}}},
		{name: "complex", args: args{value: `(?m)(?P<key>\w+):\s+(?P<value>\w+)$`}, want: StandardPattern{pattern: `(?m)(?P<key>\w+):\s+(?P<value>\w+)$`, engine: goRegexpEngine{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPattern(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStandardPattern_Matches(t *testing.T) {
	type fields struct {
		pattern string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		isRegexp2 bool
		want      bool
	}{
		{name: "regexp2 (positive): .NET-style capture groups", fields: fields{pattern: `(?<name>\d{3})`}, args: args{value: "123"}, isRegexp2: true, want: true},
		{name: "regexp2 (negative): .NET-style capture groups", fields: fields{pattern: `(?<name>\d{3})`}, args: args{value: "asdf"}, isRegexp2: true, want: false},
		{name: "regexp2 (positive): comments", fields: fields{pattern: `^(?#this is a comment)[a-zA-Z]+$`}, args: args{value: "qwertyASDF"}, isRegexp2: true, want: true},
		{name: "regexp2 (negative): comments", fields: fields{pattern: `^(?#this is a comment)[a-zA-Z]+$`}, args: args{value: "this is a comment"}, isRegexp2: true, want: false},
		{name: "regexp2 (positive): positive lookbehind with backreference", fields: fields{pattern: `((?<=keep\s)rollin') \1 \1`}, args: args{value: "keep rollin' rollin' rollin'"}, isRegexp2: true, want: true},
		{name: "regexp2 (negative): positive lookbehind with backreference", fields: fields{pattern: `((?<=keep\s)rollin') \1 \1`}, args: args{value: "keep rollin', rollin', rollin'"}, isRegexp2: true, want: false},

		{name: "regexp (positive): named ascii character class", fields: fields{pattern: `[[:alpha:]]`}, args: args{value: "This is alpha"}, isRegexp2: false, want: true},
		{name: "regexp (negative): named ascii character class", fields: fields{pattern: `[[:alpha:]]`}, args: args{value: "ꝉƕᵻ§ ᵻ§ ₪ɵꝉ Ʌ£ƿƕɅ"}, isRegexp2: false, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e regexEngine
			if tt.isRegexp2 {
				e = regexp2Engine{}
			} else {
				e = goRegexpEngine{}
			}
			s := StandardPattern{
				pattern: tt.fields.pattern,
				engine:  e,
			}
			if got := s.Matches(tt.args.value); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
