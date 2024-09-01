[IN IMPLEMENTATION] gitlab-ci-verify
===
[![GitHub Release](https://img.shields.io/github/v/release/timo-reymann/gitlab-ci-verify?label=version)](https://github.com/timo-reymann/gitlab-ci-verify/releases/latest)
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/timo-reymann/gitlab-ci-verify/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/timo-reymann/gitlab-ci-verify/tree/main)
[![codecov](https://codecov.io/gh/timo-reymann/gitlab-ci-verify/graph/badge.svg?token=4tYXDueu5D)](https://codecov.io/gh/timo-reymann/gitlab-ci-verify)
[![Renovate](https://img.shields.io/badge/renovate-enabled-green?logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAzNjkgMzY5Ij48Y2lyY2xlIGN4PSIxODkuOSIgY3k9IjE5MC4yIiByPSIxODQuNSIgZmlsbD0iI2ZmZTQyZSIgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoLTUgLTYpIi8+PHBhdGggZmlsbD0iIzhiYjViNSIgZD0iTTI1MSAyNTZsLTM4LTM4YTE3IDE3IDAgMDEwLTI0bDU2LTU2YzItMiAyLTYgMC03bC0yMC0yMWE1IDUgMCAwMC03IDBsLTEzIDEyLTktOCAxMy0xM2ExNyAxNyAwIDAxMjQgMGwyMSAyMWM3IDcgNyAxNyAwIDI0bC01NiA1N2E1IDUgMCAwMDAgN2wzOCAzOHoiLz48cGF0aCBmaWxsPSIjZDk1NjEyIiBkPSJNMzAwIDI4OGwtOCA4Yy00IDQtMTEgNC0xNiAwbC00Ni00NmMtNS01LTUtMTIgMC0xNmw4LThjNC00IDExLTQgMTUgMGw0NyA0N2M0IDQgNCAxMSAwIDE1eiIvPjxwYXRoIGZpbGw9IiMyNGJmYmUiIGQ9Ik04MSAxODVsMTgtMTggMTggMTgtMTggMTh6Ii8+PHBhdGggZmlsbD0iIzI1YzRjMyIgZD0iTTIyMCAxMDBsMjMgMjNjNCA0IDQgMTEgMCAxNkwxNDIgMjQwYy00IDQtMTEgNC0xNSAwbC0yNC0yNGMtNC00LTQtMTEgMC0xNWwxMDEtMTAxYzUtNSAxMi01IDE2IDB6Ii8+PHBhdGggZmlsbD0iIzFkZGVkZCIgZD0iTTk5IDE2N2wxOC0xOCAxOCAxOC0xOCAxOHoiLz48cGF0aCBmaWxsPSIjMDBhZmIzIiBkPSJNMjMwIDExMGwxMyAxM2M0IDQgNCAxMSAwIDE2TDE0MiAyNDBjLTQgNC0xMSA0LTE1IDBsLTEzLTEzYzQgNCAxMSA0IDE1IDBsMTAxLTEwMWM1LTUgNS0xMSAwLTE2eiIvPjxwYXRoIGZpbGw9IiMyNGJmYmUiIGQ9Ik0xMTYgMTQ5bDE4LTE4IDE4IDE4LTE4IDE4eiIvPjxwYXRoIGZpbGw9IiMxZGRlZGQiIGQ9Ik0xMzQgMTMxbDE4LTE4IDE4IDE4LTE4IDE4eiIvPjxwYXRoIGZpbGw9IiMxYmNmY2UiIGQ9Ik0xNTIgMTEzbDE4LTE4IDE4IDE4LTE4IDE4eiIvPjxwYXRoIGZpbGw9IiMyNGJmYmUiIGQ9Ik0xNzAgOTVsMTgtMTggMTggMTgtMTggMTh6Ii8+PHBhdGggZmlsbD0iIzFiY2ZjZSIgZD0iTTYzIDE2N2wxOC0xOCAxOCAxOC0xOCAxOHpNOTggMTMxbDE4LTE4IDE4IDE4LTE4IDE4eiIvPjxwYXRoIGZpbGw9IiMzNGVkZWIiIGQ9Ik0xMzQgOTVsMTgtMTggMTggMTgtMTggMTh6Ii8+PHBhdGggZmlsbD0iIzFiY2ZjZSIgZD0iTTE1MyA3OGwxOC0xOCAxOCAxOC0xOCAxOHoiLz48cGF0aCBmaWxsPSIjMzRlZGViIiBkPSJNODAgMTEzbDE4LTE3IDE4IDE3LTE4IDE4ek0xMzUgNjBsMTgtMTggMTggMTgtMTggMTh6Ii8+PHBhdGggZmlsbD0iIzk4ZWRlYiIgZD0iTTI3IDEzMWwxOC0xOCAxOCAxOC0xOCAxOHoiLz48cGF0aCBmaWxsPSIjYjUzZTAyIiBkPSJNMjg1IDI1OGw3IDdjNCA0IDQgMTEgMCAxNWwtOCA4Yy00IDQtMTEgNC0xNiAwbC02LTdjNCA1IDExIDUgMTUgMGw4LTdjNC01IDQtMTIgMC0xNnoiLz48cGF0aCBmaWxsPSIjOThlZGViIiBkPSJNODEgNzhsMTgtMTggMTggMTgtMTggMTh6Ii8+PHBhdGggZmlsbD0iIzAwYTNhMiIgZD0iTTIzNSAxMTVsOCA4YzQgNCA0IDExIDAgMTZMMTQyIDI0MGMtNCA0LTExIDQtMTUgMGwtOS05YzUgNSAxMiA1IDE2IDBsMTAxLTEwMWM0LTQgNC0xMSAwLTE1eiIvPjxwYXRoIGZpbGw9IiMzOWQ5ZDgiIGQ9Ik0yMjggMTA4bC04LThjLTQtNS0xMS01LTE2IDBMMTAzIDIwMWMtNCA0LTQgMTEgMCAxNWw4IDhjLTQtNC00LTExIDAtMTVsMTAxLTEwMWM1LTQgMTItNCAxNiAweiIvPjxwYXRoIGZpbGw9IiNhMzM5MDQiIGQ9Ik0yOTEgMjY0bDggOGM0IDQgNCAxMSAwIDE2bC04IDdjLTQgNS0xMSA1LTE1IDBsLTktOGM1IDUgMTIgNSAxNiAwbDgtOGM0LTQgNC0xMSAwLTE1eiIvPjxwYXRoIGZpbGw9IiNlYjZlMmQiIGQ9Ik0yNjAgMjMzbC00LTRjLTYtNi0xNy02LTIzIDAtNyA3LTcgMTcgMCAyNGw0IDRjLTQtNS00LTExIDAtMTZsOC04YzQtNCAxMS00IDE1IDB6Ii8+PHBhdGggZmlsbD0iIzEzYWNiZCIgZD0iTTEzNCAyNDhjLTQgMC04LTItMTEtNWwtMjMtMjNhMTYgMTYgMCAwMTAtMjNMMjAxIDk2YTE2IDE2IDAgMDEyMiAwbDI0IDI0YzYgNiA2IDE2IDAgMjJMMTQ2IDI0M2MtMyAzLTcgNS0xMiA1em03OC0xNDdsLTQgMi0xMDEgMTAxYTYgNiAwIDAwMCA5bDIzIDIzYTYgNiAwIDAwOSAwbDEwMS0xMDFhNiA2IDAgMDAwLTlsLTI0LTIzLTQtMnoiLz48cGF0aCBmaWxsPSIjYmY0NDA0IiBkPSJNMjg0IDMwNGMtNCAwLTgtMS0xMS00bC00Ny00N2MtNi02LTYtMTYgMC0yMmw4LThjNi02IDE2LTYgMjIgMGw0NyA0NmM2IDcgNiAxNyAwIDIzbC04IDhjLTMgMy03IDQtMTEgNHptLTM5LTc2Yy0xIDAtMyAwLTQgMmwtOCA3Yy0yIDMtMiA3IDAgOWw0NyA0N2E2IDYgMCAwMDkgMGw3LThjMy0yIDMtNiAwLTlsLTQ2LTQ2Yy0yLTItMy0yLTUtMnoiLz48L3N2Zz4=)](https://renovatebot.com)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=timo-reymann_gitlab-ci-verify&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=timo-reymann_gitlab-ci-verify)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=timo-reymann_gitlab-ci-verify&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=timo-reymann_gitlab-ci-verify)
[![Go Report Card](https://goreportcard.com/badge/github.com/timo-reymann/gitlab-ci-verify)](https://goreportcard.com/report/github.com/timo-reymann/gitlab-ci-verify)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=timo-reymann_gitlab-ci-verify&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=timo-reymann_gitlab-ci-verify)

<p align="center">
	<img width="300" src="https://raw.githubusercontent.com/timo-reymann/gitlab-ci-verify/main/.github/images/logo.png">
    <br />
    Validate and lint your gitlab ci files using ShellCheck, the Gitlab API and curated checks
</p>

## Features

- ShellCheck for scripts
- Validation against Pipeline Editor API for project
- Automatic detection of the current gitlab project
- Available as pre-commit hook

## Installation

### Manual

#### Linux (64-bit)

```bash
curl -LO https://github.com/timo-reymann/gitlab-ci-verify/releases/download/$(curl -Lso /dev/null -w %{url_effective} https://github.com/timo-reymann/gitlab-ci-verify/releases/latest | grep -o '[^/]*$')/gitlab-ci-verify_linux-amd64 && \
chmod +x gitlab-ci-verify_linux-amd64 && \
sudo mv gitlab-ci-verify_linux-amd64 /usr/local/bin/gitlab-ci-verify
```

#### Darwin (Intel)

##### brew

```bash
brew tap timo-reymann/gitlab-ci-verify
brew install gitlab-ci-verify
```

##### manual

```bash
curl -LO https://github.com/timo-reymann/gitlab-ci-verify/releases/download/$(curl -Lso /dev/null -w %{url_effective} https://github.com/timo-reymann/gitlab-ci-verify/releases/latest | grep -o '[^/]*$')/gitlab-ci-verify_darwin-amd64 && \
chmod +x gitlab-ci-verify_darwin-amd64 && \
sudo mv gitlab-ci-verify_darwin-amd64 /usr/local/bin/gitlab-ci-verify
```

### Install with go

```bash
go install github.com/timo-reymann/gitlab-ci-verify@latest
```

### Install with pip(x)

Using pipx you can just use the following command use gitlab-ci-verify as it is:

```sh
pipx install gitlab-ci-verify
```

If you want to use it directly using the `subprocess` module you can install it with pip:

````sh
pip install gitlab-ci-verify
````

And use the package like this:

````python
import subprocess

from gitlab_ci_verify import exec

# Run process and prefix stdout and stderr
exec.exec_with_templated_output(["--help"])

# Create a subprocess, specifying how to handle stdout, stderr
exec.create_subprocess(["--help"], stdout=subprocess.PIPE, stderr=subprocess.PIPE)

# Perform command with suppressed output and return finished proces instance,
# on that one can also check if the call was successfully
exec.exec_silently(["--version"])
````

### Supported platforms

The following platforms are supported (and have prebuilt binaries /
ready to use integration):

- Linux
    - 64-bit
    - ARM 64-bit
- Darwin
    - 64-bit
    - ARM (M1/M2)
- Windows
    - 64-bit
- pre-commit

### Where to find the latest release for your platform

#### Binaries

Binaries for all of these can be found on
the [latest release page](https://github.com/timo-reymann/gitlab-ci-verify/releases/latest).

## Usage

### Command Line

```sh
gitlab-ci-verify --help
```

### pre-commit

```yaml
- repo: https://github.com/timo-reymann/gitlab-ci-verify
  rev: main
  hooks:
    - id: gitlab-ci-verify
```

## Motivation

Unfortunately, GitLab didn't provide a tool to validate CI configuration for quite a while.
Now that changed with the `glab` CLI providing `glab ci lint` but it is quite limited and under the hood just calls the
new CI Lint API.

Throughout the years quite some tools evolved, but most of them are either outdated, painful to use or install, and
basically also provide the lint functionality from the API.

As most of the logic in pipelines is written in shell scripts via the `*script` attributes these are lacking completely
from all tools out there as well as the official lint API.

The goal of gitlab-ci-verify is to provide the stock CI Lint functionality plus shellcheck.
Completed in the future some
rules to lint that common patterns are working as intended by GitLab
and void them from being pushed and leading to unexpected behavior.

## Contributing

I love your input! I want to make contributing to this project as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the configuration
- Submitting a fix
- Proposing new features
- Becoming a maintainer

To get started please read the [Contribution Guidelines](./CONTRIBUTING.md).

## Development

### Requirements

- [Go](https://go.dev/doc/install)
- [GNU make](https://www.gnu.org/software/make/)

### Test

```sh
make test-coverage-report
```

### Build

```sh
make build
```

## Credits

This whole project wouldn't be possible with the great work of the
following libraries/tools:

- [Shellcheck by koalaman](https://github.com/koalaman/shellcheck)
- [go stdlib](https://github.com/golang/go)
- [pflag by spf13](https://github.com/spf13/pflag)
