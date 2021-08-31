# docked

A Dockerfile linter.

[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)
![Go Version](https://img.shields.io/github/go-mod/go-version/jimschubert/docked)
[![Go Build](https://github.com/jimschubert/docked/actions/workflows/build.yml/badge.svg)](https://github.com/jimschubert/docked/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimschubert/docked)](https://goreportcard.com/report/github.com/jimschubert/docked)
![Docker Pulls](https://img.shields.io/docker/pulls/jimschubert/docked)
<!-- [![codecov](https://codecov.io/gh/jimschubert/docked/branch/master/graph/badge.svg)](https://codecov.io/gh/jimschubert/docked) --> 

## Installation

### Binaries

Latest binary releases are available via [GitHub Releases](https://github.com/jimschubert/docked/releases).

### Homebrew

```
brew install jimschubert/tap/docked
```

## Usage

```shell
$ docked -h

docked is a Dockerfile linting tool which aims to pull many
best practices and recommendations from multiple sources:

  * OWASP
  * Docker Official Documentation
  * Community recommendations
  * Package manager bug trackers

Usage:
  docked [command]

Available Commands:
  analyze     Analyze a Dockerfile for issues
  completion  generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.docked.yaml)
  -h, --help            help for docked
  -v, --version         version for docked
      --viper           use Viper for configuration (default true)

Use "docked [command] --help" for more information about a command.
```

## Build

Build a local distribution for evaluation using goreleaser.

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

Build and execute locally:

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
./docked
```

## License

This project is [licensed](./LICENSE) under Apache 2.0.
