# Authentication with GitLab

The tool takes a few sources into consideration in the following order when authenticating with GitLab:

- `--gitlab-token` commandline flag
- [.netrc](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html) in your home folder.
  For example:

  ```text
  machine gitlab.com password <api_token>
  ```

  The *default* machine is not used.
- Vault token specified via `--gitlab-token vault://<path>#<field>` with environment variable `VAULT_ADDR` set to base
  url for
  vault, and either `VAULT_TOKEN` set or `~/.vault-token` present
- `GITLAB_TOKEN` environment variable

For the project detection, all git remote URLs of the repository are tried, and the first URL that returns a valid API
response is used. In case you cloned via SSH it tries to convert it to the HTTPs host automatically. If the ssh URL
differs from the HTTPs url you should specify it manually using the `--gitlab-base-url`, without protocol e.g.
`--gitlab-base-url git.example.com`

## Using the GitLab CI Lint API Proxy

The `--gitlab-base-url` can also be set to a GitLab CI Lint API Proxy instance. This is useful when you want to
centralize GitLab API access or run the linting in environments where direct GitLab access is restricted.

The proxy exposes the same API path as the regular GitLab API (`/api/v4/projects/{project}/ci/lint`) and forwards lint
requests to the upstream GitLab instance. This makes it a drop-in replacement for direct GitLab API access.

Note that the proxy only handles the CI lint endpoint - it is not a full GitLab API proxy.

The proxy uses the exact same credential resolution logic as described [above](#authentication-with-gitlab).

### Configuration

| Option | Type | Description |
|--------|------|-------------|
| `--gitlab-base-url` | Flag (required) | Base URL of the upstream GitLab instance |
| `GITLAB_TOKEN` | Environment variable | GitLab access token (or use other auth methods) |
| `LOG_LEVEL` | Environment variable | Set to `verbose` or `debug` for logging, silent by default |

The proxy always listens on port `8080`.

### Example Docker Compose setup

```yaml
services:
  ci-lint-proxy:
    image: timoreymann/gitlab-ci-verify/gitlab-ci-lint-api-proxy
    command:
      - --gitlab-base-url=gitlab.com
    environment:
      - GITLAB_TOKEN=${GITLAB_TOKEN}
    ports:
      - "8080:8080"
```

You can then point `gitlab-ci-verify` to the proxy:

```bash
gitlab-ci-verify --gitlab-base-url=localhost:8080
```
