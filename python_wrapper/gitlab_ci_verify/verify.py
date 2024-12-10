from os import PathLike

from gitlab_ci_verify.config import GitlabCiVerifyConfig
from gitlab_ci_verify.parse import _parse_output
from gitlab_ci_verify.process import _execute
from typing import Unpack

def verify_file(
        root: PathLike | str,
        file: str | None = None,
        **configuration: Unpack[GitlabCiVerifyConfig],
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
        **configuration: Unpack[GitlabCiVerifyConfig],
):
    """
    Verify gitlab ci file contents using gitlab-ci-verify

    :param root: Root folder of repository
    :param content: CI YAML contents
    :param configuration: Configuration to use for verification
    """
    proc = _execute(root, file="-", stdin=content, **configuration)
    return _parse_output(proc)
