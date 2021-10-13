package rules

import (
	"sort"
	"strings"

	"github.com/jimschubert/docked/model"
	"github.com/jimschubert/docked/model/docker/commands"
	"github.com/jimschubert/docked/model/shell"
	"github.com/jimschubert/docked/model/validations"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

func sortInstallerArgs() validations.Rule {
	commandLookup := installIndicators()

	r := validations.MultiContextRule{
		Name:             "sort-installer-args",
		Summary:          "Sort installed packages for package managers: apt-get, apk, npm, etc.",
		Details:          "Sorting installed packages alphabetically prevents duplicates and simplifies maintainability.",
		Priority:         model.LowPriority,
		Commands:         []commands.DockerCommand{commands.Run},
		URL:              model.StringPtr("https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#sort-multi-line-arguments"),
		AppliesToBuilder: true,
		Evaluator: validations.MultiContextPerNodeEvaluator{
			Fn: func(node *parser.Node, validationContext validations.ValidationContext) model.Valid {
				posixCommands, err := shell.NewPosixCommandFromNode(node)
				if err != nil {
					log.Warnf("Unable to parse RUN command, validation not evaluated: %#v", node.Location())
					return model.Skipped
				}

				managers := commandLookup.Keys()
				for _, command := range posixCommands {
					var name string
					var argIndexStart = 0
					name = strings.TrimLeft(command.Name, `\`)
					// this is a naive "best-guess" means to support finding package manager in some edge-cases
					if name == "sudo" || name == "su" || name == "gosu" {
						for idx, arg := range command.Args {
							if !strings.HasPrefix(arg, "-") {
								name = strings.TrimLeft(arg, `\`)
								argIndexStart = idx
								break
							}
						}
					}

					if model.StringSliceContains(&managers, name) {
						// We assume all commands are format:
						// package-manager [options] <command> [<args>...]
						// we need to find the command, then evaluate the args
						var seenInstallCommand bool
						packages := make([]string, 0)
						for _, arg := range command.Args[argIndexStart:] {
							if !strings.HasPrefix(arg, "-") {
								if !seenInstallCommand {
									seenInstallCommand = commandLookup[name](arg)
									continue
								}
								packages = append(packages, arg)
							}
						}

						if !sort.SliceIsSorted(packages, func(i, j int) bool {
							return packages[i] < packages[j]
						}) {
							return model.Recommendation
						}
					}
				}

				return model.Success
			},
		},
	}
	return &r
}

func installIndicators() model.PredicateMap {
	commandLookup := model.PredicateMap{
		"apt": func(s string) bool {
			// see https://manpages.ubuntu.com/manpages/xenial/man8/apt.8.html
			return s == "install"
		},
		"apt-get": func(s string) bool {
			// see https://linux.die.net/man/8/apt-get
			return s == "install"
		},
		"yum": func(s string) bool {
			// see https://man7.org/linux/man-pages/man8/yum.8.html
			return s == "install"
		},
		"apk": func(s string) bool {
			// see https://wiki.alpinelinux.org/wiki/Alpine_Linux_package_management#Add_a_Package
			// NOTE: apk switches can come _after_ packages
			return s == "add"
		},
		"npm": func(s string) bool {
			// see https://docs.npmjs.com/cli/v7/commands/npm-install
			return s == "install" || s == "i" || s == "add" || s == "isntall"
		},
	}
	return commandLookup
}

func init() {
	AddRule(sortInstallerArgs())
}
