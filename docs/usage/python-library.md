Automate validation using the python library
===

If you want to automate checking GitLab CI YAML files for errors, you can use the python library.

It uses the gitlab-ci-verify binary under the hood, which is distributed using
the [gitlab-ci-verify-bin](https://pypi.org/project/gitlab-ci-verify-bin) package.

Use cases for the python library include:

- dynamically generated CI pipelines
- bulk check of CI YAML files across repositories
- bulk check of CI templates

## Code sample

```python
from gitlab_ci_verify import verify_file

# Verify .gitlab-ci.yml in /path/to/repo is valid
valid, findings = verify_file("/path/to/repo")

# verify include.yml in /path/to/repo is valid
valid, findings = verify_file("/path/to/repo", "include.yml")

# or if you want to verify file content for a given repository
# valid, findings = verify_content("/path/to/repo","ci-yaml content")

print(f"Valid:    {valid}")
print(f"Findings: {findings}")
```

## API documentation

See the [python wrapper documentation](https://gitlab-ci-verify.timo-reymann.de/python-wrapper/) for 
comprehensive API documentation.