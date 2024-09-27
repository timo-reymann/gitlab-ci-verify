import subprocess
from os import PathLike
from typing import Unpack

from gitlab_ci_verify_bin.exec import create_subprocess

from gitlab_ci_verify.config import GitlabCiVerifyConfig


def _add_arg_if_set(args: list[str], config: Unpack[GitlabCiVerifyConfig], key: str, flag: str):
    val = config.get(key, None)
    if val is not None:
        args.extend(
            [
                flag,
                val,
            ]
        )


def _execute(
        root: PathLike | str,
        file: str | None = None,
        stdin=None,
        **config: Unpack[GitlabCiVerifyConfig]
):
    args = [
        "--format",
        "json"
    ]

    _add_arg_if_set(args, config, "gitlab_base_url", "--gitlab-base-url")
    _add_arg_if_set(args, config, "gitlab_token", "--gitlab-token")
    _add_arg_if_set(args, config, "fail_severity", "--severity")

    if file is not None:
        args.extend(
            [
                "--gitlab-ci-file",
                file,
            ]
        )

    excluded_checks = config.get("excluded_checks", None)
    if excluded_checks is not None:
        for check in excluded_checks:
            args.extend(
                [
                    "--exclude",
                    check,
                ]
            )

    proc = create_subprocess(
        args,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        stdin=subprocess.PIPE if stdin is not None else None,
        cwd=root,
    )

    if stdin is not None:
        proc.stdin.write(stdin)
        proc.stdin.close()

    proc.wait()
    return proc
