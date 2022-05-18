package common

import (
	"github.com/DG9Jww/gatherInfo/logger"
	"github.com/panjf2000/ants/v2"
)

func NewPool(num int) *ants.Pool {
	pool, err := ants.NewPool(num)
	if err != nil {
		logger.ConsoleLog(logger.ERROR, err.Error())
	}
	return pool
}
