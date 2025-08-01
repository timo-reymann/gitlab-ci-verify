How shellcheck is executed
===

[shellcheck](https://www.shellcheck.net/) is a linter to find bugs in shell scripts.

shellcheck is used to evaluate and validate `script` blocks in GitLab CI YAML files, including:

- `before_script`
- `script`
- `after_script`

## Distribution of shellcheck

To avoid you having to install shellcheck, it comes pre-bundled with gitlab-ci-verify.

The version output also gives you information about the used shellcheck version:

```sh
gitlab-ci-verify --version
```

In addition, you can find information about the latest version in the source tree
at [internal/shellcheck/bin/README.txt](https://github.com/timo-reymann/gitlab-ci-verify/blob/main/internal/shellcheck/bin/README.txt).

## Execution of the binary

Under the hood gitlab-ci-verify uses [memexec](https://github.com/amenzhinsky/go-memexec). Depending on the platform, the
execution method differs:

- macOS: copy to a temporary file and execute
- Windows: copy to a temporary file and execute
- Linux: Creates a memory file descriptor that loads the entire binary in the RAM and executes it

For efficiency reasons, the shellcheck instance is reused for all scripts discovered.