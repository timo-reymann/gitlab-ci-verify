from dataclasses import dataclass, fields


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
        class_fields = {f.name for f in fields(cls)}
        return Finding(**{k: v for k, v in dict_val.items() if k in class_fields})
