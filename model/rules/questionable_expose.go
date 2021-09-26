package rules

import (
	"bytes"
	"fmt"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/docker/types"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/sirupsen/logrus"
)

var (
	commonPorts = types.ExposeList{
		types.Expose{PortRange: singlePort(21), Description: model.StringPtr("FTP")},
		types.Expose{PortRange: singlePort(22), Description: model.StringPtr("SSH")},
		types.Expose{PortRange: singlePort(23), Description: model.StringPtr("Telnet")},
		types.Expose{PortRange: singlePort(25), Description: model.StringPtr("SMTP")},
		types.Expose{PortRange: singlePort(43), Description: model.StringPtr("WHOIS")},
		types.Expose{PortRange: singlePort(53), Description: model.StringPtr("DNS")},
		types.Expose{PortRange: types.PortRange{Start: 67, End: 68}, Description: model.StringPtr("DHCP")},
		types.Expose{PortRange: singlePort(69), Description: model.StringPtr("TFTP")},
		types.Expose{PortRange: singlePort(110), Description: model.StringPtr("POP3")},
		types.Expose{PortRange: singlePort(115), Description: model.StringPtr("SFTP")},
		types.Expose{PortRange: singlePort(143), Description: model.StringPtr("IMAP")},
		types.Expose{PortRange: singlePort(873), Description: model.StringPtr("rsync")},

		types.Expose{PortRange: singlePort(993), Description: model.StringPtr("IMAP SSL")},
		types.Expose{PortRange: singlePort(955), Description: model.StringPtr("POP3 SSL")},
		types.Expose{PortRange: singlePort(1080), Description: model.StringPtr("SOCKS")},
		types.Expose{PortRange: singlePort(3128), Description: model.StringPtr("Proxy List")},
		types.Expose{PortRange: singlePort(3306), Description: model.StringPtr("MySQL")},
		types.Expose{PortRange: singlePort(3389), Description: model.StringPtr("RDP")},
		types.Expose{PortRange: singlePort(5432), Description: model.StringPtr("PostgreSQL")},
		types.Expose{PortRange: singlePort(5900), Description: model.StringPtr("VNC")},
		types.Expose{PortRange: singlePort(5938), Description: model.StringPtr("TeamViewer")},
	}
)

func singlePort(input int) types.PortRange {
	return types.PortRange{Start: input, End: input}
}

func questionableExpose() validations.Rule {
	portDescriptionFormatted := func(e types.Expose) string {
		buf := bytes.Buffer{}
		buf.WriteString(fmt.Sprintf("%d", e.PortRange.Start))
		if e.PortRange.End > e.PortRange.Start {
			buf.WriteString(fmt.Sprintf("-%d", e.PortRange.End))
		}
		if e.Protocol != "" {
			buf.WriteString(fmt.Sprintf("/%s", e.Protocol))
		}
		if e.Description != nil {
			buf.WriteString(fmt.Sprintf(" (%s)", *e.Description))
		}
		return buf.String()
	}

	r := validations.MultiContextRule{
		Name:     "questionable-expose",
		Summary:  "Avoid documenting EXPOSE with sensitive ports",
		Details:  "The EXPOSE command is metadata and does not actually open ports. Documenting the intention to expose sensitive ports poses a security concern.",
		Commands: []commands.DockerCommand{commands.Expose},
		Evaluator: validations.MultiContextPerNodeEvaluator{
			Fn: func(node *parser.Node, validationContext validations.ValidationContext) model.Valid {
				trimStart := len(node.Value) + 1 // command plus trailing space
				defs := node.Original[trimStart:]
				exposeList, err := types.ParseExposeList(defs)
				if err != nil {
					logrus.WithError(err).Debugf("Unable to parse list of exposed ports at line %d.", validationContext.Locations[0].Start.Line)
					return model.Failure
				}

				questionable := false
				for _, expose := range exposeList {
					for _, c := range commonPorts {
						if expose.PortRange.Intersects(c.PortRange) {
							logrus.Infof("Found questionable port exposed %s.", portDescriptionFormatted(c))
							questionable = true
						}
					}
				}

				if questionable {
					return model.Failure
				}
				return model.Success
			},
		},
	}
	return &r
}

func init() {
	AddRule(questionableExpose())
}
