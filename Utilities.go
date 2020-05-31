package pcrawler

import (
	"io"
	"os"
)

//func AIDtoURL(aid string) (string, error) {
//	return "", nil
//}
//
//func URLToAID(url string) (string, error) {
//	return "", nil
//}

func ErrorStripper(result interface{}, err error) interface{} {
	if err != nil {
		return nil
	}
	return result
}

func GetResourceFromFile(filePath string) io.ReadCloser {
	ret, _ := os.Open(filePath)
	return ret
}
