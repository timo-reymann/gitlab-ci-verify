# Installation

gitlab-ci-verify supports a wide variety of platforms.

## pre-commit

To check your Gitlab CI YAML before pushing or even as part of CI itself.

```yaml
- repo: https://github.com/timo-reymann/gitlab-ci-verify
  rev: v2.6.0
  hooks:
    - id: gitlab-ci-verify
```

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

