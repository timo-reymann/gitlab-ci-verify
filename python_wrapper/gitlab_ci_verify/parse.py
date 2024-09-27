import json
from json import JSONDecodeError

from gitlab_ci_verify import Finding

class FindingsParseException(Exception):
    pass

def _parse_output(proc):
    valid = proc.returncode == 0
    stdout_content = proc.stdout.read()
    try:
        raw_findings = json.loads(stdout_content)
    except JSONDecodeError:
        raise FindingsParseException(proc.stderr.read())
    parsed_findings = []
    for raw_finding in raw_findings:
        parsed_findings.append(Finding.from_dict(raw_finding))

    return valid, parsed_findings