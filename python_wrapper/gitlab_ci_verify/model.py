from dataclasses import dataclass


@dataclass(frozen=True, kw_only=True)
class Finding:
    severity: str
    code: str
    line: int
    message: str
    link: str
    file: str

    @classmethod
    def from_dict(cls, dict_val: dict[str, any]):
        return Finding(**dict_val)
