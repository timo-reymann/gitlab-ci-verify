import json

from gitlab_ci_verify import Finding


def _parse_output(proc):
    valid = proc.returncode == 0
    raw_findings = json.loads(proc.stdout.read())
    parsed_findings = []
    for raw_finding in raw_findings:
        parsed_findings.append(Finding.from_dict(raw_finding))

    return valid, parsed_findings