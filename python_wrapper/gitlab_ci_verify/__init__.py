from os import PathLike

from gitlab_ci_verify.model import Finding
from gitlab_ci_verify.parse import _parse_output
from gitlab_ci_verify.process import _execute


def verify(
        root: PathLike | str,
        gitlab_base_url: str | None = None,
        gitlab_ci_file: str | None = None,
        gitlab_token: str | None = None,
        excluded_checks: list[str] = None,
        fail_severity: str | None = None,
):
    """
    Verify gitlab ci file using gitlab-ci-verify

    :param root: Root folder of repository
    :param gitlab_base_url: Override gitlab base url
    :param gitlab_ci_file: Override gitlab ci file
    :param gitlab_token: Specify explicit gitlab token
    :param excluded_checks: Exclude set of checks
    :param fail_severity: On which severity findings should be considered an error
    :return: Validity of pipeline as first argument and findings as second
    """
    proc = _execute(root, gitlab_base_url, gitlab_ci_file, gitlab_token, excluded_checks, fail_severity)
    return _parse_output(proc)
