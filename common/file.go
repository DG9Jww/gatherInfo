package common

import (
	"os"

	"github.com/DG9Jww/gatherInfo/logger"
)

func LoadFile(path string) *os.File {
	handle, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
	}
	return handle
}

func WriteFile(path string) *os.File {
	handle, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
	}
	return handle
}
