# Local setup for gitlab-ci-verify development

If you want to work on gitlab-ci-verify you can do so, easily.

Before starting to develop make sure you have read
the [Contribution Guidelines](https://github.com/timo-reymann/gitlab-ci-verify/blob/main/CONTRIBUTING.md).

## Requirements

Make sure to install the following tools:

- [Go](https://go.dev/doc/install)
- [npm](https://www.npmjs.com/get-npm)
- [GNU make](https://www.gnu.org/software/make/)
- [Python 3.10+](https://www.python.org/downloads/)

## Run 

```sh
go run main.go -h
```

## Test

Run tests and show the coverage report in your default browser.

```sh
make test-coverage-report
```

## Build

Build the binary.

```sh
make build
```

