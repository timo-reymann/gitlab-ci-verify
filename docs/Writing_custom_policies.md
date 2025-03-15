Writing custom policies
===

You can write custom policies for your projects
using [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/).

To get started to create the folder `.gitlab-ci-verify/checks` in your project and add a file with the extension
`.rego`.

```rego
# package does not really matter, as long as is it does **not** contain gitlab_ci_verify
package my_project_checks

# Import helpers
import data.gitlab_ci_verify

# Use the latest and greatest v1 API
import rego.v1

# Define a rule that ensures each image has a tag
findings contains gitlab_ci_verify.error("PROJ-1001", sprintf("Job %s does not contain tag for image", [job]), yamlPathToLineNumber(sprintf(".%s.image", [job]))) if {
    some job in input.yaml[job]
    not contains(input.yaml[job].image, ":")
}

# Define a rule that nesures each pipeline has the job pages

findings contains finding if {
  not input.yaml.pages
  
  finding := gitlab_ci_verify.error("PROJ-1002", "The pipeline needs to contain a pages job", -1)
}
```

## Provided helpers

Helpers are methods you can import and that are written in Rego to make it easier to write policies.

| Signature                                                       | Description                                                   |
|:----------------------------------------------------------------|:--------------------------------------------------------------|
| `gitlab_ci_verify.error(code, message, line)`                   | Create a finding with the given code, message and line        |
| `gitlab_ci_verify.error_with_link(code, message, line, link)`   | Create a finding with the given code, message, line and link  |
| `gitlab_ci_verify.warning(code, message, line)`                 | Create a warning with the given code, message and line        |
| `gitlab_ci_verify.warning_with_link(code, message, line, link)` | Create a warning with the given code, message, line and link  |
| `gitlab_ci_verify.info(code, message, line)`                    | Create an info with the given code, message and line          |
| `gitlab_ci_verify.info_with_link(code, message, line, link)`    | Create an info with the given code, message, line and link    |
| `gitlab_ci_verify.create_finding(level, code, message, line)`   | Create a finding with the given level, code, message and line |

## Extra builtins

The following additional builtins are available in the Rego policies:

| Signature                    | Description                                                                                                                                                                                                              |
|:-----------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `yamlPathToLineNumber(path)` | Convert a [YAML path](https://github.com/vmware-labs/yaml-jsonpath?tab=readme-ov-file#syntax) to a line number, if nothing is found returns `-1` which is accepted for findings, when they can’t be correlated to a line |

## Input

The input for the policies is based on the parsed gitlab ci file.

```json5
{
  "yaml": {
    // ... the parsed YAML from the disk
  },
  "mergedYaml": {
    // ... the merged YAML from the API, it is set to nil when run with the `--no-lint-api-in-ci` flag
  }
}
```

## Remote bundles

You can use remote bundles to share policies between projects.
To do so you can use the `--include-opa-bundle` flag to specify a bundle to include.

The bundle should be a tarball containing the rego files, built
using [opa](https://www.openpolicyagent.org/docs/latest/cli/#opa-build).

To allow caching, the server should support [RFC7232](https://datatracker.ietf.org/doc/html/rfc7232) compliant caching
headers.

## Further resources

- [Rego playground](https://play.openpolicyagent.org/)
- [Library for gitlab-ci-verify check helpers](./internal/rego_policies/lib.rego)