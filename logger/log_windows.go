package logger

import (
	"log"
	"os"
	"syscall"
)

type logType struct {
	color      int
	prefix     string
	isProgress bool
	row        int
	colu       int
	total      int
}

//windows colour
const (
	BLACK = iota
	BLUE
	GREEN
	CYAN
	RED
	PURPLE
	YELLOW
	LIGHT_GRAY
	GRAY
	LIGHT_BLUE
	LIGHT_GREEN
	LIGHT_CYAN
	LIGHT_RED
	LIGHT_PURPLE
	LIGHT_YELLOW
	WHITE
)

var (
	// loggerTime  *log.Logger       = log.New(os.Stdout, " ", log.Ldate|log.Ltime)
	kernel32    *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	proc        *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
	CloseHandle *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)
)

func ConsoleLog(t logType, v string) {
	logger := log.New(os.Stdout, "", 0)
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(t.color))
	logger.Println(t.prefix, v)
	defer CloseHandle.Call(handle)
}

func ConsoleLog2(t logType, v ...interface{}) {
	logger := log.New(os.Stdout, "", 0)
	//logger.Print("[", t.color, t.prefix, ENDC, "]", v)
	logger.Print("[")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(t.color))
	logger.Print(t.prefix)
	defer CloseHandle.Call(handle)
	logger.Print("]")
	logger.Print(v)
}

func LogToFile(file *os.File, v ...interface{}) {
	logger := log.New(file, "", 0)
	logger.Print(v...)
}
