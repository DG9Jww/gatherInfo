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

//func ConsoleLog(t logType, v ...interface{}) {
//	logger * log.Logger = log.New(os.Stdout, "", 0)
//
//	if t.isProgress {
//		//Move Curson to Specific Position
//		logger.Printf("\033[%d;%df", t.row, t.colu)
//		handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(t.color))
//		logger.Println(t.prefix, v)
//	}
//	//If not progress bar
//	fmt.Printf("\033[u")
//	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(t.color))
//	logger.Println(t.prefix, v)
//	fmt.Printf("\033[s")
//	defer CloseHandle.Call(handle)
//}

func ConsoleLog(t logType, v ...interface{}) {
	logger := log.New(os.Stdout, "", 0)
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(t.color))
	logger.Println(t.prefix, v)
	defer CloseHandle.Call(handle)
}

func LogToFile(file *os.File, v ...interface{}) {
	logger := log.New(file, "", 0)
	logger.Print(v...)
}
