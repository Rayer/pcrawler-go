package PttUtils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
)

type IndexInfo struct {

}

type IndexInitialParameters struct {
	MaxIndex int
	LastDocUrl string
	PinnedDocs int
}

type DocumentMetrics struct {

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
	headIndexUrl := c.CreateHeadIndexUrl()
	res, err := c.GetHttpResponse(headIndexUrl)
	if err != nil {
		return IndexInitialParameters{}, err
	}

	// Required information :
	// Index count, Last timestamp, Last document url hash
	ret := IndexInitialParameters{}

	//dump, err := ioutil.ReadAll(contentReader)
	//fmt.Println(string(dump))

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return IndexInitialParameters{}, err
	}

	//Index Count
	doc.Find("div#action-bar-container a").Each(func(i int, selection *goquery.Selection) {
		link, exists := selection.Attr("href")
		if !exists {
			return
		}
		r := regexp.MustCompile("/bbs/.+/index(\\d+)\\.html")
		if find := r.FindStringSubmatch(link); len(find) > 1 {
			index, _ := strconv.Atoi(find[1])
			if index + 1 > ret.MaxIndex {
				ret.MaxIndex = index + 1
			}
		}
	})

	ret.PinnedDocs = len(doc.Find(".r-list-sep").NextAll().Nodes)

	var lastDoc *goquery.Selection
	if ret.PinnedDocs > 0 {
		lastDoc = doc.Find(".r-list-sep").Prev()
	} else {
		lastDoc = doc.Find(".r-ent").Last()
	}

	lastDocUrl, _ := lastDoc.Find(".title a").Attr("href")
	lastDocUrl = "https://www.ptt.cc" + lastDocUrl
	ret.LastDocUrl = lastDocUrl

	return ret, nil
}

func (c *Crawler) CreateHeadIndexUrl() *url.URL {
	if c.Board == "" {
		panic("Not initialized, use NewCrawler() to create initialized Crawler!")
	}

	if ret, err := url.Parse(fmt.Sprintf("https://www.ptt.cc/bbs/%s/index.html", c.Board)); err != nil {
		panic("Internal error...")
	} else {
		return ret
	}
}

func (c *Crawler) GetContent(url *url.URL) (string, error) {

	reader, err := c.GetHttpResponse(url)
	if err != nil {
		return "", err
	}

	ba, err := ioutil.ReadAll(reader.Body)
	if err != nil {
		return "", err
	}

	return string(ba), nil
}

func (c* Crawler) GetHttpResponse(url *url.URL) (*http.Response, error) {
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