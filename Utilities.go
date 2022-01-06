package pcrawler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
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

func CreatePttClient(u *url.URL) http.Client {
	var cookies []*http.Cookie
	cookie := &http.Cookie{Name: "over18", Value: "1", Domain: "ptt.cc", Path: "/"}
	cookies = append(cookies, cookie)
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, cookies)
	return http.Client{Timeout: 3 * time.Second, Jar: jar}
}
