package logger

var (
	STATUS = logType{color: WHITE, prefix: "[STATUS]"}
	INFO   = logType{color: LIGHT_BLUE, prefix: "[INFO]"}
	WARN   = logType{color: LIGHT_YELLOW, prefix: "[WARN]"}
	ERROR  = logType{color: LIGHT_RED, prefix: "[ERROR]"}
	NORMAL = logType{color: WHITE, prefix: "[+]"}

	//dirscan
	R30X = logType{color: LIGHT_YELLOW}
	R20X = logType{color: BLUE}
	R40X = logType{color: LIGHT_RED}
	R50X = logType{color: PURPLE}

	//subdomain
	SUBDOMAIN = logType{color: LIGHT_GREEN, prefix: "[✓]"}

	//portscan
	PORTSCAN = logType{color: LIGHT_BLUE}

	//fingerprint
	FINGERPRINT = logType{color: LIGHT_YELLOW}
)
