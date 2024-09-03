package checks

var SeverityError = 0
var SeverityWarning = 1
var SeverityInfo = 2
var SeverityStyle = 3

type CheckFinding struct {
	Severity int
	Code     string
	FileName string
	Line     int
	Message  string
	Link     string
}

func (cf *CheckFinding) SeverityName() string {
	switch cf.Severity {
	case SeverityError:
		return "Error"
	case SeverityWarning:
		return "Warning"
	case SeverityInfo:
		return "Info"
	case SeverityStyle:
		return "Style"
	}
	return ""
}
