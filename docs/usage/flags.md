# CLI Flags

gitlab-ci-verify's behavior can be controlled with the following command-line-flags:

| Flag | Shorthand | Default | Description |
| ---- | --------- | ------- | ----------- |
| gitlab-ci-file | | .gitlab-ci.yml | The Yaml file used to configure GitLab CI |
| gitlab-base-url | | | Set the gitlab base url explicitly in case detection does not work or your clone and base url differs. |
| gitlab-token | | | Gitlab token to use, if not specified the netrc is evaluated and if that also does not contain credentials, tries to load the environment variable GITLAB_TOKEN. |
| shellcheck-flags | | | Pass custom flags to shellcheck |
| format | f | text | Format for the output, valid options are json, table, text and gitlab. If GITLAB_CI_VERIFY_OUTPUT_FORMAT this parameter is ignored. |
| output | o | | Write the report to a file instead of stdout. |
| severity | S | style | Set the severity level on which to consider findings as errors and exiting with non zero exit code. |
| exclude | E | | Exclude the given check codes. |


| Flag | Shorthand | Description |
| ---- | --------- | ----------- |
| debug | | Enable debug output |
| verbose | | Enable verbose output |
| offline | | Add this flag to avoid validating against Pipeline Check API, so only the offline checks are performed. Please note that checks relying on the merged YAML will also not be executed in that case. |
| no-lint-api-in-ci | | Add this flag to avoid validating against Pipeline Check API, as its assumed that running in CI is proof enough the synt is valid. Please note that checks relying on the merged YAML will also not be executed in that case. |
| include-opa-bundle | I | Include remote OPA bundles for checks |
