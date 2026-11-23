# CLI Flags

gitlab-ci-verify's behavior can be controlled with various command-line-flags.

You can get a full list by running

```bash
gitlab-ci-verify --help
```

## Validating a feature branch (`--lint-ref`)

By default the CI Lint API resolves `include` entries, such as CI/CD components, against the default branch
of your project. When your pipeline includes a component pinned to the current commit, for example
`component: $CI_SERVER_FQDN/$CI_PROJECT_PATH/ci@$CI_COMMIT_SHA`, and that component only exists on your
feature branch, the validation fails or uses the stale component from the default branch.

Pass `--lint-ref` with a branch or tag to validate the pipeline in the context of the ref.
In this mode the CI Lint API runs a pipeline simulation (dry run) instead of the static check:

```bash
gitlab-ci-verify --lint-ref "$CI_COMMIT_REF_NAME"
```
