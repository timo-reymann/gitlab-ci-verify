import unittest
from unittest.mock import MagicMock

from gitlab_ci_verify.parse import FindingsParseException, _parse_output


class ParseTest(unittest.TestCase):
    def test_empty(self):
        proc = MagicMock()
        proc.stdout.read.return_value = "[]"
        proc.returncode = 0
        valid, findings = _parse_output(proc)
        assert valid
        assert len(findings) == 0

    def test_findings(self):
        proc = MagicMock()
        proc.stdout.read.return_value = '[{"severity": "STYLE", "code": "CH-123","line": 1, "message": "test", "link": "https://check.example.com", "file": "file.yml"}]'
        proc.returncode = 0
        valid, findings = _parse_output(proc)
        assert valid
        assert len(findings) == 1

    def test_err(self):
        proc = MagicMock()
        proc.stdout.read.return_value = ""
        proc.stderr.read.return_value = "error message"
        proc.returncode = 1
        with self.assertRaises(FindingsParseException):
            _parse_output(proc)
