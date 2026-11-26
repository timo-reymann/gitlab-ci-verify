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

## ShellCheck for the merged configuration (`--shellcheck-merged-yaml`)

By default, ShellCheck only analyzes scripts from your local files: the entry file and its local includes.
This keeps every finding mapped to a file and line you can fix.

Scripts that only exist after GitLab resolved the configuration aren't covered by that:

- jobs from CI/CD components, project, remote and template includes
- scripts with `$[[ inputs.* ]]` interpolations, which GitLab resolves

Pass `--shellcheck-merged-yaml` to also analyze the `merged_yaml` returned by the CI Lint API:

```bash
gitlab-ci-verify --shellcheck-merged-yaml
```

Scripts that are identical to a local script, and lines that are unchanged compared to the local job, are skipped,
so findings are not reported twice.

The merged YAML contains no information about where a script came from. Findings from it are therefore marked with
`merged:` in the message, and include the affected script line. They're reported for the job in your local files
if it exists there, otherwise for the entry file without a line number.

This requires the CI Lint API. When it isn't used (`--no-lint-api`, `--no-lint-api-in-ci`) or the pipeline is
invalid, only local scripts are analyzed.
