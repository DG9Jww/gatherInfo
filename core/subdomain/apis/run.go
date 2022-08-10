package apis

import (
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/DG9Jww/gatherInfo/logger"
)

var resSlice []Result

type Result struct {
	domain string
	ip     string
}

func Run(domains []string) []Result {

	//some API need special process
	//go func() {
	//	for _, d := range domains {
	//		err, s := scripts.StartCrt(d)
	//		if err != nil {
	//			logger.ConsoleLog(logger.ERROR, err.Error())
	//		}
	//		for _, v := range s {
	//			res := Result{}
	//			res.domain = v
	//			resSlice = append(resSlice, res)
	//		}
	//	}
	//}()

	// range script directory
	var wg sync.WaitGroup
	fileSystem := os.DirFS(rootDir)
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.ConsoleLog(logger.ERROR, err.Error())
		}
		if !strings.HasSuffix(path, "json") {
			return nil
		}

		APIName := strings.Split(path, ".")[0]
		data, err := fs.ReadFile(fileSystem, path)
		for _, d := range domains {
			wg.Add(1)
			go start(APIName, data, d, &wg)
		}
		return nil
	})
	wg.Wait()

	for _, v := range resSlice {
		logger.ConsoleLog2(logger.CustomizeLog(logger.BLUE, v.domain), v.ip)
	}
	return resSlice
}

func init() {
	//functions
	funcList[funcBase64] = myBase64
}
