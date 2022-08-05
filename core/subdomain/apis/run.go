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
		wg.Wait()
		return nil
	})
	return resSlice
}

func init() {
	//functions
	funcList[funcBase64] = myBase64
}
