package docker

import (
	"strings"
	"testing"

	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func TestInstruction(t *testing.T) {
	type args struct {
		node *parser.Node
	}

	a := func(input string) args {
		cmd := strings.FieldsFunc(input, func(r rune) bool {
			return r == ' '
		})[0]
		return args{node: &parser.Node{Original: input, Value: strings.ToLower(cmd)}}
	}

	tests := []struct {
		name     string
		args     args
		wantCmd  commands.DockerCommand
		wantText string
	}{
		{name: "FROM", args: a("FROM scratch"), wantCmd: commands.From, wantText: "scratch"},
		{name: "WORKDIR", args: a("WORKDIR /go/src/app"), wantCmd: commands.Workdir, wantText: "/go/src/app"},
		{name: "ADD", args: a("ADD . /go/src/app"), wantCmd: commands.Add, wantText: ". /go/src/app"},
		{name: "RUN", args: a("RUN apk --no-cache add gcc g++ make ca-certificates && apk add git"), wantCmd: commands.Run, wantText: "apk --no-cache add gcc g++ make ca-certificates && apk add git"},
		{name: "COPY", args: a("COPY --from=builder /copy/copy /copy"), wantCmd: commands.Copy, wantText: "--from=builder /copy/copy /copy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotText := Instruction(tt.args.node)
			if gotCmd != tt.wantCmd {
				t.Errorf("Instruction() gotCmd = %v, wantCmd %v", gotCmd, tt.wantCmd)
			}
			if gotText != tt.wantText {
				t.Errorf("Instruction() gotText = %v, wantText %v", gotText, tt.wantText)
			}
		})
	}
}
