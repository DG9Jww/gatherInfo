package logger

import (
	"log"
	"os"
)

type logType struct {
	color  string
	prefix string
}

const (
	// Color
	WHITE        = "\033[37m"
	GREEN        = "\033[92m"
	LIGHT_GREEN  = "\033[1;92m"
	RED          = "\033[91m"
	LIGHT_RED    = "\033[1;91m"
	YELLOW       = "\033[33m"
	LIGHT_YELLOW = "\033[93m"
	BLUE         = "\033[94m"
	LIGHT_BLUE   = "\033[1;94m"
	LIGHT_WHITE  = "\033[97m"
	PURPLE       = "\033[35m"
	ENDC         = "\033[0m"
)

func ConsoleLog(t logType, v string) {
	logger := log.New(os.Stdout, "", 0)
	logger.Print(t.color, t.prefix, v, ENDC)
}

func LogToFile(file *os.File, v ...interface{}) {
	logger := log.New(file, "", 0)
	logger.Print(v...)
}

func ConsoleLog2(t logType, v string) {
	logger := log.New(os.Stdout, "", 0)
	logger.Printf("[%s%s%s] %s", t.color, t.prefix, ENDC, v)
}
