# docked

A Dockerfile linting tool which aims to pull many best practices and recommendations from multiple sources:

* OWASP
* Docker Official Documentation
* Community recommendations
* Package manager bug trackers

[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)
![Go Version](https://img.shields.io/github/go-mod/go-version/jimschubert/docked)
[![Go Build](https://github.com/jimschubert/docked/actions/workflows/build.yml/badge.svg)](https://github.com/jimschubert/docked/actions/workflows/build.yml)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/jimschubert/docked?color=orange&label=Docker%20Image%20Size)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimschubert/docked)](https://goreportcard.com/report/github.com/jimschubert/docked)
<!-- [![codecov](https://codecov.io/gh/jimschubert/docked/branch/master/graph/badge.svg)](https://codecov.io/gh/jimschubert/docked) --> 

## tldr;

```
docked analyze ./Dockerfile
```

Outputs:
![](./.github/screens/output.png)

And, it's customizable (you can ignore, re-prioritize, or add custom rules via regex).

## Install

### Binaries

Latest binary releases are available via [GitHub Releases](https://github.com/jimschubert/docked/releases).

### Homebrew

```
brew install jimschubert/tap/docked
```

### Docker

```
docker pull jimschubert/docked:latest
```

When running the docker image, be sure to mount and reference the sources appropriately. For example:

### Completions

After you've installed the binary either manually or via Homebrew, consider enabling completions for your shell. 

For instructions, view help for your target shell.

#### zsh

```
docked completion zsh --help
```

#### bash

```
docked completion bash --help
```

#### fish

```
docked completion fish --help
```

#### powershell

```
docked completion powershell --help
```

## Usage

```shell
$ docked analyze --help

Analyze a Dockerfile for issues
If not provided, FILE defaults to ./Dockerfile

Usage:
  docked analyze [FILE] [flags]

Flags:
  -h, --help                   help for analyze
  -i, --ignore strings         The lint ids to ignore
  -k, --no-buildkit-warnings   Whether to suppress Docker parser warnings
      --regex-engine string    The regex engine to use (regexp, regexp2) (default "regexp2")
      --report-type string     The type of reporting output (text, json, html) (default "text")

Global Flags:
      --config string   config file (default is $HOME/.docked.yaml)
      --viper           use Viper for configuration (default true)
```

Things to consider:

* Buildkit warnings should be disabled when piping output (for example when using `--report-type json`), but this is _not forced_
* The `regexp2` engine is default because it supports full regular expression syntax. Compare differences in [regexp2's README](https://github.com/dlclark/regexp2#compare-regexp-and-regexp2). Note that `regexp2` patterns are not run in compatibility mode in docked, although that might change later.
* `viper` configuration is work-in-progress. Feel free to contribute.

## Configuration

The optional configuration file follows this example syntax:

```
ignore:
  - D7:tagged-latest
rule_overrides:
  'D5:secret-aws-access-key': low
custom_rules:
  - name: custom-name
    summary: Your custom summary
    details: Your additional rule details
    pattern: '.' # some regex pattern
    priority: critical
    command: add
```

## Build

Build a local distribution for evaluation using goreleaser (easiest).

```bash
goreleaser release --skip-publish --snapshot --rm-dist
```

This will create an executable application for your os/architecture under `dist`:

```
dist
├── docked_darwin_amd64
│   └── docked
├── docked_linux_386
│   └── docked
├── docked_linux_amd64
│   └── docked
├── docked_linux_arm64
│   └── docked
├── docked_linux_arm_6
│   └── docked
└── docked_windows_amd64
    └── docked.exe
```

Build and execute locally using go:

* Get dependencies

```shell
go get -d ./...
```

* Build

```shell
go build -o docked ./cmd/docked/
```

* Run

```shell
./docked --help
```

## License

This project is [licensed](./LICENSE) under Apache 2.0.
