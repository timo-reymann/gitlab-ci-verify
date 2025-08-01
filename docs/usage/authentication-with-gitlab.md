# Authentication with GitLab

The tool takes a few sources into consideration in the following order when authenticating with GitLab:

- `--gitlab-token` commandline flag
- [netrc](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html) in your home folder
- Vault token specified via `--gitlab-token vault://<path>#<field>` with environment variable `VAULT_ADDR` set to base
  url for
  vault, and either `VAULT_TOKEN` set or `~/.vault-token` present
- `GITLAB_TOKEN` environment variable

For the project detection, all git remote URLs of the repository are tried, and the first URL that returns a valid API
response is used. In case you cloned via SSH it tries to convert it to the HTTPs host automatically. If the ssh URL
differs from the HTTPs url you should specify it manually using the `--gitlab-base-url`, without protocol e.g.
`--gitlab-base-url git.example.com`
