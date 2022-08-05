package apis

import "encoding/base64"

func myBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
