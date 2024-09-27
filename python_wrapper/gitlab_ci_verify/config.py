import typing


class GitlabCiVerifyConfig(typing.TypedDict):
    gitlab_base_url: typing.Optional[str]
    """Override gitlab base url"""
    gitlab_token: typing.Optional[str]
    """Specify explicit gitlab token"""
    excluded_checks: typing.Optional[list[str]]
    """Exclude set of checks"""
    fail_severity: typing.Optional[str]
    """On which severity findings should be considered an error"""
