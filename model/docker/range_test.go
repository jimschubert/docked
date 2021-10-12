package docker

import (
	"reflect"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func ppos(line int, character int) parser.Position {
	return parser.Position{Line: line, Character: character}
}

func pos(line int, character int) Position {
	return Position{Line: line, Character: character}
}

func TestFromParserRanges(t *testing.T) {
	type args struct {
		p []parser.Range
	}
	tests := []struct {
		name string
		args args
		want []Location
	}{
		{
			name: "Single range",
			args: args{p: []parser.Range{
				{Start: ppos(11, 2), End: ppos(12, 24)},
			}},
			want: []Location{
				{Start: pos(11, 2), End: pos(12, 24)},
			},
		},
		{
			name: "multiple ranges",
			args: args{p: []parser.Range{
				{Start: ppos(1, 20), End: ppos(12, 24)},
				{Start: ppos(5, 2), End: ppos(200, 123)},
				{Start: ppos(0, 0), End: ppos(1, 1)},
				{Start: ppos(110, 110), End: ppos(110, 110)},
			}},
			want: []Location{
				{Start: pos(1, 20), End: pos(12, 24)},
				{Start: pos(5, 2), End: pos(200, 123)},
				{Start: pos(0, 0), End: pos(1, 1)},
				{Start: pos(110, 110), End: pos(110, 110)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromParserRanges(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromParserRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
