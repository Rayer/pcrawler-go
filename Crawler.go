package PttUtils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type IndexInfo struct {

}

type IndexInitialParameters struct {
	MaxIndex int
	StartDate time.Time
	PinnedDocs int
	LastDocumentUrlHash uint32
}


type Crawler struct {
	Board string
	IndexInfo []IndexInfo
	sharedClient http.Client
}

func NewCrawler(boardName string) (*Crawler, error) {
	ret := Crawler{Board:boardName}
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Jar: cookieJar,
	}
	ret.sharedClient = client
	return &ret, nil
}

func (c *Crawler) GetIndexInitialParameters() (IndexInitialParameters, error) {
	headIndexUrl := c.createHeadIndexUrl()
	var contentReader io.Reader
	if con, err := c.getContentReader(headIndexUrl); err != nil {
		return IndexInitialParameters{}, err
	} else {
		contentReader = con
	}

	// Required information :
	// Index count, Last timestamp, Last document url hash
	ret := IndexInitialParameters{}

	doc, err := goquery.NewDocumentFromReader(contentReader)
	if err != nil {
		return IndexInitialParameters{}, err
	}



	return ret, nil
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

	reader, err := c.getContentReader(url)
	if err != nil {
		return "", err
	}

	ba, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(ba), nil
}

func (c* Crawler) getContentReader(url *url.URL) (io.Reader, error) {
	var cookies []*http.Cookie
	cookie := &http.Cookie{Name: "over18", Value: "1", Domain: "ptt.cc", Path: "/"}
	cookies = append(cookies, cookie)
	c.sharedClient.Jar.SetCookies(url, cookies)

	res, err := c.sharedClient.Get(url.String())
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return nil, err
	}

	return res.Body, nil
}