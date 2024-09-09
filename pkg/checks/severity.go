package checks

import "strings"

func SeverityLevelToName(level int) string {
	switch level {
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

func SeverityNameToLevel(name string) int {
	switch strings.ToLower(name) {
	case "error":
		return SeverityError
	case "warning":
		return SeverityWarning
	case "info":
		return SeverityInfo
	case "style":
		return SeverityStyle
	}
	return -1
}
