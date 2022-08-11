package apis

import (
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/DG9Jww/gatherInfo/common"
	"github.com/DG9Jww/gatherInfo/logger"
)


const (
	//function
	funcBase64 = `$base64`
)


var resSlice []string
var lock sync.Mutex
var (
	funcList = make(map[string]func(string) string)

	//
	rootDir = `core/subdomain/apis/scripts`
)

//return subdomain slice
func Run(domains []string) []string {

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

	//remove duplicates
	resSlice = common.RemoveStringDuplicate(resSlice)
	return resSlice
}

func addResSlice(item string) {
	lock.Lock()
	resSlice = append(resSlice, item)
	lock.Unlock()
}

func init() {
	//functions
	funcList[funcBase64] = myBase64
}
