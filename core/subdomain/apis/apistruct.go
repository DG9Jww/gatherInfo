package apis

import (
	"fmt"
	_ "reflect"
)

//var APIStruct = map[string]reflect.Type {
//	"virustotal": reflect.TypeOf(virustotal{}),
//}

var APIStruct = map[string]interface{}{
	"virustotal": virustotal{},
}

//virustotal
type virustotal struct {
	Data []virustotal2 `json:"data"`
}

type virustotal2 struct {
	ID string `json:"id"`
}

func xxx(i interface{}) {
	switch i.(type) {
	case virustotal:
		fmt.Println("convert to  virustotal")
	default:
		fmt.Println("convert failed")
	}

}
