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
- Curated checks for common mistakes (feel free to [contribute new ones](https://gitlab-ci-verify.timo-reymann.de/add-builtin-check.html))
- Automatic detection of the current gitlab project with an option to overwrite
- Available as pre-commit hook
- Usable to valid dynamically generated pipelines using the [python wrapper](https://gitlab-ci-verify.timo-reymann.de/usage/python-library.html)
- Support for *gitlab.com* and self-hosted instances
- Support for [custom policies](https://gitlab-ci-verify.timo-reymann.de/extending/writing-custom-policies.html) written
  in [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/)
- Resolve and validate
  includes ([how it works and limitations](https://gitlab-ci-verify.timo-reymann.de/how-it-works/include-resolution.html))

## Installation

See the [Installation section](https://gitlab-ci-verify.timo-reymann.de/installation.html) in the documentation.

## Documentation

You can find the full documentation on [GitHub Pages](https://gitlab-ci-verify.timo-reymann.de/), including:

- How it works
- How to add new checks
- How to write custom policies using rego
- How to authenticate with GitLab

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

To get started, please read the [Contribution Guidelines](./CONTRIBUTING.md).

## Credits

This whole project wouldn't be possible with the great work of the
following libraries/tools:

- [Shellcheck by koalaman](https://github.com/koalaman/shellcheck)
- [go stdlib](https://github.com/golang/go)
- [pflag by spf13](https://github.com/spf13/pflag)
- [go-yaml](https://github.com/go-yaml/yaml), which I forked
  to [timo-reymann/go-yaml](https://github.com/timo-reymann/go-yaml)
