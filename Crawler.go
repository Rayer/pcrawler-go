package pcrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func NewCrawler(boardName string) *Crawler {
	ret := Crawler{Board:boardName}
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

func (c *Crawler) GetIndexInitialParameters() (IndexInitialParameters, error) {
	headIndexUrl := c.createHeadIndexUrl()
	res, err := c.getHttpResponse(headIndexUrl)
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

func (c *Crawler) ParseDocument(url *url.URL) (*PDocRaw, error) {

	res, err := c.getHttpResponse(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		fmt.Printf("Error : %s\n", err)
		return nil, err
	}

	var infoList []CommitInfo

	doc.Find(".push").Each(func(i int, s *goquery.Selection) {
		//fmt.Println(s.Text())
		name := s.Find(".push-userid").Text()
		ctime := strings.TrimSpace(strings.Replace(s.Find(".push-ipdatetime").Text(), "\n", "", -1))
		ctype := s.Find(".push-tag").Text()
		content := s.Find(".push-content").Text()

		//Take last 11 of time string, because sometimes we will got IP as prefix....
		var ctimeTransformed time.Time
		if len(ctime) >= 11 {
			ctimeTransformed, err = time.Parse("01/02 15:04", ctime[len(ctime) - 11 : ])
			if err != nil {
				fmt.Printf("Error parsing %s, set to NOW() \n", ctime)
				ctimeTransformed = time.Now()
			}
		} else {
			fmt.Printf("Error parsing %s, set to NOW() \n", ctime)
			ctimeTransformed = time.Now()
		}

		var ctypeTransformed int
		switch ctype {
		case "推 ":
			ctypeTransformed = 1
		case "噓 ":
			ctypeTransformed = -1
		case "→ ":
			ctypeTransformed = 0
		default:
			fmt.Printf("Can't transform ctype : %s\n", ctype)
		}

		//fmt.Printf("Find %s(%s)(%s) : %s\n", name, ctime, ctype, content)
		infoList = append(infoList, CommitInfo{
			Type:      ctypeTransformed,
			Committer: name,
			Timestamp: ctimeTransformed,
			Content:   content,
		})
	})

	/*
		Here is an example for commit log
		<div class="push"><span class="f1 hl push-tag">→ </span><span class="f3 hl push-userid">deann</span><span class="f3 push-content">: 台灣哪間公司法遵能力強的 笑死 不就存乎老闆一心而已</span><span class="push-ipdatetime"> 05/28 09:52
		   </span>
		</div>
	*/

	//raw, _ := doc.Html()
	ret := &PDocRaw{
		Title:             doc.Find("div#main-content").Find(":nth-child(3)").Find(".article-meta-value").Text(),
		Author:            doc.Find("div#main-content").Find(":nth-child(1)").Find(".article-meta-value").Text(),
		//RawArticleHtml:    raw,
		PublicUrl:         url.String(),
		CommitterInfoList: infoList,
		ProcessTime:       time.Now(),
	}

	//fmt.Printf("Prcessed : %s with Committer info list length of : %d", ret.Title, len(ret.CommitterInfoList))

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

func (c* Crawler) getHttpResponse(url *url.URL) (*http.Response, error) {
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