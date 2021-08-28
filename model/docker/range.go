package docker

import (
	"bytes"
	"fmt"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Location is the start and end Position of a docker instruction.
type Location struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position is the line number and starting character in some block of text (a Dockerfile).
type Position struct {
	Line      int `json:"line,omitempty"`
	Character int `json:"character,omitempty"`
}

func (r Location) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%d:%d", r.Start.Line, r.Start.Character))
	if r.Start.Line != r.End.Line && r.Start.Character != r.End.Character {
		buf.WriteString(fmt.Sprintf("-%d:%d", r.End.Line, r.End.Character))
	}
	return buf.String()
}

func fromRange(p parser.Range) Location {
	return Location{
		Start: Position{
			Line:      p.Start.Line,
			Character: p.Start.Character,
		},
		End: Position{
			Line:      p.End.Line,
			Character: p.End.Character,
		},
	}
}

// FromParserRanges converts Docker sdk parser.Range into a tool specific Location
func FromParserRanges(p []parser.Range) []Location {
	ranges := make([]Location, 0)
	for _, parserRange := range p {
		ranges = append(ranges, fromRange(parserRange))
	}
	return ranges
}
