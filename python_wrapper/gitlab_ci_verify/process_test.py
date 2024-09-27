import subprocess
import unittest
from unittest.mock import MagicMock, patch

from gitlab_ci_verify.process import _execute


class ProcessTest(unittest.TestCase):
    @patch("gitlab_ci_verify.process.create_subprocess")
    def test_execute_no_args(self, create_subprocess_mock: MagicMock):
        _execute("/workspace/folder")
        create_subprocess_mock.assert_called_with(
            ['--format', 'json'],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            stdin=None,
            cwd='/workspace/folder'
        )

    @patch("gitlab_ci_verify.process.create_subprocess")
    def test_execute_all_args(self, create_subprocess_mock: MagicMock):
        _execute(
            "/workspace/folder",
            file="test.yml",
            excluded_checks=["CHECK-123"],
            fail_severity="warning",
            gitlab_token="token",
            gitlab_base_url="git.company.org",
            stdin="blub",
        )
        create_subprocess_mock.assert_called_with(
            [
                '--format', 'json',
                '--gitlab-base-url', 'git.company.org',
                '--gitlab-token', 'token',
                '--severity', 'warning',
                '--gitlab-ci-file', 'test.yml',
                '--exclude', 'CHECK-123'
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            stdin=subprocess.PIPE,
            cwd='/workspace/folder'
        )
