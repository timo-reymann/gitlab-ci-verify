import typing


class GitlabCiVerifyConfig(typing.TypedDict):
    gitlab_base_url: typing.NotRequired[str]
    """Override gitlab base url"""
    gitlab_token: typing.NotRequired[str]
    """Specify explicit gitlab token"""
    excluded_checks: typing.NotRequired[list[str]]
    """Exclude set of checks"""
    fail_severity: typing.NotRequired[str]
    """On which severity findings should be considered an error"""
