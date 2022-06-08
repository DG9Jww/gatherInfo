package logger

import "fmt"

var (
	INFO   = logType{color: LIGHT_BLUE, prefix: fmt.Sprintf("[%s%s%s]", BLUE, "INFO", ENDC)}
	WARN   = logType{color: LIGHT_YELLOW, prefix: fmt.Sprintf("[%s%s%s]", YELLOW, "WARN", ENDC)}
	ERROR  = logType{color: LIGHT_RED, prefix: fmt.Sprintf("[%s%s%s]", LIGHT_RED, "ERROR", ENDC)}
	NORMAL = logType{color: WHITE, prefix: fmt.Sprintf("[%s%s%s]", WHITE, "+", ENDC)}

	//dirscan
	R30X = logType{color: LIGHT_YELLOW}
	R20X = logType{color: BLUE}
	R40X = logType{color: LIGHT_RED}
	R50X = logType{color: PURPLE}

	//subdomain
	SUBDOMAIN  = logType{color: LIGHT_GREEN, prefix: "[✓]"}
	SUBDOMAIN2 = logType{color: LIGHT_BLUE, prefix: "[✓]"}

	//portscan
	PORTSCAN = logType{color: LIGHT_BLUE}

	//fingerprint
	FINGERPRINT = logType{color: LIGHT_YELLOW}
)

func CustomizeLog(color string, prefix string) logType {
	return logType{color: color, prefix: prefix}
}
