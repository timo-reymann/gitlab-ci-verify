# Ignoring findings

You can ignore findings by adding comments in the format `# gitlab-ci-verify: ignore:<check_id>` to your CI YAML files.

This works in several places:

## In the same line as the finding:

Place the ignore comment at the end of the affected line

  ```yaml
  pages:
    artifacts: { }# gitlab-ci-verify: ignore:GL-201
  ```

## In the line above the finding:

Place the ignore comment above the line with the finding

  ```yaml
  pages:
    # gitlab-ci-verify: ignore:GL-201
    artifacts: { }
  ```

## Globally for the file of the finding:

Add a ignore comment to the top of the file to ignore a finding for an entire file.

  ```yaml
  # gitlab-ci-verify: ignore:GL-201
  pages:
    artifacts: { }
  ```
