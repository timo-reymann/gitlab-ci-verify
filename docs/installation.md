# Installation

gitlab-ci-verify supports a wide variety of platforms.

## pre-commit

To check your Gitlab CI YAML before pushing or even as part of CI itself.

```yaml
- repo: https://github.com/timo-reymann/gitlab-ci-verify
  rev: v2.10.0
  hooks:
    - id: gitlab-ci-verify
```

## GitLab CI Template

For easy integration into your GitLab CI/CD pipelines, you can use the provided CI template. This template includes a pre-configured job that runs gitlab-ci-verify and generates code quality reports.

### Basic Usage

1. Include the template in your `.gitlab-ci.yml`:

    ```yaml
    include:
      - template: 'ci-templates/v2.gitlab-ci.yml'
    ```

2. Set the required GitLab token variable:

    ```yaml
    lint-ci:
      variables:
        GITLAB_CI_VERIFY_GITLAB_TOKEN: "$YOUR_GITLAB_ACCESS_TOKEN"
    ```

### Advanced Configuration

You can customize the template by overriding variables:

```yaml
include:
  - remote: 'https://gitlab-ci-verify.timo-reymann.de/ci-templates/v2.gitlab-ci.yml'

lint-ci:
  variables:
    GITLAB_TOKEN: "$YOUR_GITLAB_ACCESS_TOKEN"
    GITLAB_CI_VERIFY_SEVERITY: "warning"  # Change severity level
    GITLAB_CI_VERIFY_CI_YAML: "custom-ci.yml"  # Verify different CI file
    GITLAB_CI_VERIFY_EXTRA_ARGS: "--exclude PROJ-1001"  # Exclude specific checks
```

### Features

- **Automatic Validation**: Runs on every merge request
- **Code Quality Integration**: Shows findings directly in GitLab's MR interface
- **Customizable**: Adjust severity levels, exclude checks, and more
- **Debugging**: Provides raw JSON output for troubleshooting

For more details, see the [inline documentation in the template](https://github.com/timo-reymann/gitlab-ci-verify/blob/main/ci-templates/v2.gitlab-ci.yml).

## Containerized

If you prefer to use containerized workflows, use the provided OCI image.

```sh
docker run --rm -it -v $PWD:/workspace -e GITLAB_TOKEN="your token" timoreymann/gitlab-ci-verify
```

## Install with pipx

Using pipx you can just use the following command use gitlab-ci-verify as it is:

```sh
pipx install gitlab-ci-verify-bin
```

## Install as library using pip

If you want to automate validation of pipelines using pipelines, you can use
the [gitlab-ci-verify python package](https://pypi.org/project/gitlab-ci-verify/).

````sh
pip install gitlab-ci-verify
````

For more details check the [Use with python guide](./usage/python-library.md)

## Manual

### Linux (64-bit)

```bash
curl -LO https://github.com/timo-reymann/gitlab-ci-verify/releases/download/$(curl -Lso /dev/null -w %{url_effective} https://github.com/timo-reymann/gitlab-ci-verify/releases/latest | grep -o '[^/]*$')/gitlab-ci-verify_linux-amd64 && \
chmod +x gitlab-ci-verify_linux-amd64 && \
sudo mv gitlab-ci-verify_linux-amd64 /usr/local/bin/gitlab-ci-verify
```

### Darwin (Intel)

```bash
curl -LO https://github.com/timo-reymann/gitlab-ci-verify/releases/download/$(curl -Lso /dev/null -w %{url_effective} https://github.com/timo-reymann/gitlab-ci-verify/releases/latest | grep -o '[^/]*$')/gitlab-ci-verify_darwin-amd64 && \
chmod +x gitlab-ci-verify_darwin-amd64 && \
sudo mv gitlab-ci-verify_darwin-amd64 /usr/local/bin/gitlab-ci-verify
```

### Windows

Download the latest [release](https://github.com/timo-reymann/gitlab-ci-verify/releases/latest) for Windows and put in
your `PATH`.

## Where to find the latest release for your platform

### Binaries

Binaries for all of these can be found on
the [latest release page](https://github.com/timo-reymann/gitlab-ci-verify/releases/latest).

### Docker

For the docker image, check the [docker hub](https://hub.docker.com/r/timoreymann/gitlab-ci-verify).

