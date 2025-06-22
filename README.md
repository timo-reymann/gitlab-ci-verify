gitlab-ci-verify
===
[![GitHub Release](https://img.shields.io/github/v/release/timo-reymann/gitlab-ci-verify?label=version)](https://github.com/timo-reymann/gitlab-ci-verify/releases/latest)
[![PyPI version](https://badge.fury.io/py/gitlab-ci-verify.svg)](https://pypi.org/project/gitlab-ci-verify)
[![PyPI - Downloads](https://img.shields.io/pypi/dm/gitlab-ci-verify)](https://pypi.org/project/gitlab-ci-verify)
[![DockerHub Pulls](https://img.shields.io/docker/pulls/timoreymann/gitlab-ci-verify)](https://hub.docker.com/r/timoreymann/gitlab-ci-verify)
[![GitHub all releases download count](https://img.shields.io/github/downloads/timo-reymann/gitlab-ci-verify/total)](https://github.com/timo-reymann/gitlab-ci-verify/releases)
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
    Validate and lint your gitlab ci files using ShellCheck, the Gitlab API, curated checks or even build your own checks
</p>

## Features

- ShellCheck for scripts
- Validation against Pipeline Lint API for project
- Curated checks for common mistakes (feel free to [contribute new ones](./docs/checks/Add_check.md))
- Automatic detection of the current gitlab project with an option to overwrite
- Available as pre-commit hook
- Usable to valid dynamically generated pipelines using the [python wrapper](#install-as-library-using-pip)
- Support for *gitlab.com* and self-hosted instances
- Support for [custom policies](#writing-custom-policies) written
  in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/)
- Resolve and validate includes ([more information on how it works and limitations](docs/Include_resolution.md))

### Example output

| Format |                                                                Screenshot                                                                |
|:-------|:----------------------------------------------------------------------------------------------------------------------------------------:|
| text   |  ![Text output screenshot](https://raw.githubusercontent.com/timo-reymann/gitlab-ci-verify/main/.github/images/example_output/text.png)  |
| json   |  ![JSON output screenshot](https://raw.githubusercontent.com/timo-reymann/gitlab-ci-verify/main/.github/images/example_output/json.png)  |
| table  | ![Table output screenshot](https://raw.githubusercontent.com/timo-reymann/gitlab-ci-verify/main/.github/images/example_output/table.png) |

## Installation

### [pre-commit](#as-pre-commit-hook)

### [docker](#containerized)

### Install with pipx

Using pipx you can just use the following command use gitlab-ci-verify as it is:

```sh
pipx install gitlab-ci-verify-bin
```

### Install as library using pip

If you want to use it directly using the `subprocess` module you can install it with pip:

````sh
pip install gitlab-ci-verify
````

And use the package like this:

````python
from gitlab_ci_verify import verify_file

# Verify .gitlab-ci.yml in /path/to/repo is valid
valid, findings = verify_file("/path/to/repo")

# verify include.yml in /path/to/repo is valid
valid, findings = verify_file("/path/to/repo", "include.yml")

# or if you want to verify file content for a given repository
# valid, findings = verify_content("/path/to/repo","ci-yaml content")

print(f"Valid:    {valid}")
print(f"Findings: {findings}")
````

Also see the [python wrapper documentation](https://timo-reymann.github.io/gitlab-ci-verify/python-wrapper/)

### Manual

#### Linux (64-bit)

```bash
curl -LO https://github.com/timo-reymann/gitlab-ci-verify/releases/download/$(curl -Lso /dev/null -w %{url_effective} https://github.com/timo-reymann/gitlab-ci-verify/releases/latest | grep -o '[^/]*$')/gitlab-ci-verify_linux-amd64 && \
chmod +x gitlab-ci-verify_linux-amd64 && \
sudo mv gitlab-ci-verify_linux-amd64 /usr/local/bin/gitlab-ci-verify
```

#### Darwin (Intel)

```bash
curl -LO https://github.com/timo-reymann/gitlab-ci-verify/releases/download/$(curl -Lso /dev/null -w %{url_effective} https://github.com/timo-reymann/gitlab-ci-verify/releases/latest | grep -o '[^/]*$')/gitlab-ci-verify_darwin-amd64 && \
chmod +x gitlab-ci-verify_darwin-amd64 && \
sudo mv gitlab-ci-verify_darwin-amd64 /usr/local/bin/gitlab-ci-verify
```

#### Windows

Download the latest [release](https://github.com/timo-reymann/gitlab-ci-verify/releases/latest) for Windows and put in
your `PATH`.

### Install with go

```bash
go install github.com/timo-reymann/gitlab-ci-verify@latest
```

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
- pre-commit (x86 & ARM)
- Docker (x86 & ARM)

### Where to find the latest release for your platform

#### Binaries

Binaries for all of these can be found on
the [latest release page](https://github.com/timo-reymann/gitlab-ci-verify/releases/latest).

#### Docker

For the docker image, check the [docker hub](https://hub.docker.com/r/timoreymann/gitlab-ci-verify).

## Usage

### Command Line

```sh
gitlab-ci-verify --help
```

### Containerized

```sh
docker run --rm -it -v $PWD:/workspace -e GITLAB_TOKEN="your token" timoreymann/gitlab-ci-verify
```

### As pre-commit hook

```yaml
- repo: https://github.com/timo-reymann/gitlab-ci-verify
  rev: v2.1.11
  hooks:
    - id: gitlab-ci-verify
```

## Authentication with GitLab

The tool takes a few sources into consideration in the following order when authenticating with GitLab:

- `--gitlab-token` commandline flag
- [netrc](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html) in your home folder
- Vault token specified via `--gitlab-token vault://<path>#<field>` with environment variable `VAULT_ADDR` set to base
  url for
  vault, and either `VAULT_TOKEN` set or `~/.vault-token` present
- `GITLAB_TOKEN` environment variable

For the project detection, all git remote URLs of the repository are tried, and the first URL that returns a valid API
response is used. In case you cloned via SSH it tries to convert it to the HTTPs host automatically. If the ssh URL
differs from the HTTPs url you should specify it manually using the `--gitlab-base-url`, without protocol e.g.
`--gitlab-base-url git.example.com`

## Ignoring findings

You can ignore findings by adding comments in the format `# gitlab-ci-verify: ignore:<check_id>` to your CI YAML files.

This works in several places:

- In the same line as the finding:
  ```yaml
  pages:
    artifacts: {}# gitlab-ci-verify: ignore:GL-201
  ```
- In the line above the finding:
  ```yaml
  pages:
    # gitlab-ci-verify: ignore:GL-201
    artifacts: {}
  ```
- Globally for the file of the finding:
  ```yaml
  # gitlab-ci-verify: ignore:GL-201
  pages:
    artifacts: {}
  ```

## Writing custom policies

You can write custom policies for your projects
using [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/).

Find more information in the dedicated [documentation](./docs/Writing_custom_policies.md) for custom policies.

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
- [npm](https://www.npmjs.com/get-npm)
- [GNU make](https://www.gnu.org/software/make/)
- [Python 3.10+](https://www.python.org/downloads/)

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
- [go-yaml](https://github.com/go-yaml/yaml), which I forked to [timo-reymann/go-yaml](https://github.com/timo-reymann/go-yaml)
