import subprocess
from os import PathLike

from gitlab_ci_verify_bin.exec import create_subprocess


def execute(
        root: PathLike | str,
        gitlab_base_url: str | None = None,
        gitlab_ci_file: str | None = None,
        gitlab_token: str | None = None,
        excluded_checks: list[str] = None,
        fail_severity: str | None = None,
):
    args = [
        "--format",
        "json"
    ]

    if gitlab_base_url is not None:
        args.extend(
            [
                "--gitlab-base-url",
                gitlab_base_url,
            ]
        )

    if gitlab_ci_file is not None:
        args.extend(
            [
                "--gitlab-ci-file",
                gitlab_ci_file,
            ]
        )

    if gitlab_token is not None:
        args.extend(
            [
                "--gitlab-token",
                gitlab_token,
            ]
        )

    if excluded_checks is not None:
        for check in excluded_checks:
            args.extend(
                [
                    "--exclude",
                    check,
                ]
            )
    if fail_severity is not None:
        args.extend(
            [
                "--severity",
                fail_severity
            ]
        )

    proc = create_subprocess(
        args,
        subprocess.PIPE,
        subprocess.PIPE,
        cwd=root,
    )
    proc.wait()
    return proc
