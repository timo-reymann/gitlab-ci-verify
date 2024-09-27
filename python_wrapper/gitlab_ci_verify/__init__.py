from os import PathLike

from gitlab_ci_verify.config import GitlabCiVerifyConfig
from gitlab_ci_verify.model import Finding
from gitlab_ci_verify.parse import _parse_output
from gitlab_ci_verify.process import _execute


def verify_file(
        root: PathLike | str,
        file: str | None = None,
        **configuration: GitlabCiVerifyConfig,
):
    """
    Verify gitlab ci file using gitlab-ci-verify

    :param root: Root folder of repository
    :param file: Override gitlab ci file to validate
    :param configuration: Configuration to use for verification
    """
    proc = _execute(root, file=file, **configuration)
    return _parse_output(proc)


def verify_content(
        root: PathLike | str,
        content: str,
        **configuration: GitlabCiVerifyConfig,
):
    """
    Verify gitlab ci file contents using gitlab-ci-verify

    :param root: Root folder of repository
    :param content: CI YAML contents
    :param configuration: Configuration to use for verification
    """
    proc = _execute(root, file="-", stdin=content, **configuration)
    return _parse_output(proc)


if __name__ == "__main__":
    # valid, findings = verify_file("/Users/phpe/workspace/backend-csharp")
    valid, findings = verify_content("/Users/phpe/workspace/backend-csharp", "{}")
    print(valid, findings)
