package pcrawler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Crawler struct {
	Board        string
	sharedClient http.Client
}

func NewCrawler(boardName string) *Crawler {
	ret := Crawler{Board: boardName}
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic("Error creating cookie jar")
	}
	client := http.Client{
		Jar: cookieJar,
	}
	ret.sharedClient = client
	return &ret
}

func (c *Crawler) GetIndexInitialParameters() (*IndexOverview, error) {
	headIndexUrl := c.createHeadIndexUrl()
	res, err := c.getHttpResponse(headIndexUrl)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			logrus.Warnf("Fail to close connection : %s", err.Error())
		}
	}()
	return ParseIndexOverview(res.Body)
}

func (c *Crawler) ParseDocument(url *url.URL) (*PDocRaw, error) {
	return ParseSingleRawDocument(url.String())
}

func (c *Crawler) ParseIndex(index int) {
	//indexUrl := c.createIndexUrl()
}

func (c *Crawler) createHeadIndexUrl() *url.URL {
	if c.Board == "" {
		panic("Not initialized, use NewCrawler() to create initialized Crawler!")
	}

	if ret, err := url.Parse(fmt.Sprintf("https://www.ptt.cc/bbs/%s/index.html", c.Board)); err != nil {
		panic("Internal error...")
	} else {
		return ret
	}
}

func (c *Crawler) getContent(url *url.URL) (string, error) {

	reader, err := c.getHttpResponse(url)
	if err != nil {
		return "", err
	}

	ba, err := ioutil.ReadAll(reader.Body)
	if err != nil {
		return "", err
	}

	return string(ba), nil
}

func (c *Crawler) getHttpResponse(url *url.URL) (*http.Response, error) {
	var cookies []*http.Cookie
	cookie := &http.Cookie{Name: "over18", Value: "1", Domain: "ptt.cc", Path: "/"}
	cookies = append(cookies, cookie)
	c.sharedClient.Jar.SetCookies(url, cookies)

	res, err := c.sharedClient.Get(url.String())
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return nil, err
	}

	return res, nil
}
