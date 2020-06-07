package pcrawler

import (
	"encoding/json"
	"io"
	"io/ioutil"
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

func DumpExpectedResult(object interface{}, filename string) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, bytes, 0644)
	return err
}
