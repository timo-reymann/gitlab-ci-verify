from unittest import TestCase

from gitlab_ci_verify import Finding


class ModelTest(TestCase):
    def test_from_dict_valid(self):
        finding = Finding.from_dict(
            {
                "severity": "info",
                "code": "GL-XXX",
                "line": 1,
                "message": "Test",
                "link": "https://link.example",
                "file": "test.yml",
            }
        )
        assert finding == Finding(
            severity='info',
            code='GL-XXX',
            line=1,
            message='Test',
            link='https://link.example',
            file='test.yml'
        )

    def test_from_dict_extra_fields(self):
        finding = Finding.from_dict(
            {
                "severity": "info",
                "code": "GL-XXX",
                "line": 1,
                "message": "Test",
                "link": "https://link.example",
                "file": "test.yml",
                "extraField": "blub",
            }
        )
        assert finding == Finding(
            severity='info',
            code='GL-XXX',
            line=1,
            message='Test',
            link='https://link.example',
            file='test.yml'
        )

    def test_from_dict_less_fields(self):
        with self.assertRaises(TypeError):
            Finding.from_dict(
                {
                    "severity": "info",
                    "code": "GL-XXX",
                    "line": 1,
                    "message": "Test",
                    "link": "https://link.example",
                }
            )
