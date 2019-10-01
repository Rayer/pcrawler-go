package PttUtils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func CreateRawDocument(fromUrl string) (*PDocRaw, error) {

	client, _ := createCrawlerClient(fromUrl)

	res, err := client.Get(fromUrl)

	if err != nil {
		fmt.Printf("Error : %s\n", err)
	}
	//body, err := ioutil.ReadAll(res.Body)
	//fmt.Printf("Body : %s", body)

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
		var ctime_transformed time.Time
		if len(ctime) >= 11 {
			ctime_transformed, err = time.Parse("01/02 15:04", ctime[len(ctime) - 11 : ])
			if err != nil {
				fmt.Printf("Error parsing %s, set to NOW() \n", ctime)
				ctime_transformed = time.Now()
			}
		} else {
			fmt.Printf("Error parsing %s, set to NOW() \n", ctime)
			ctime_transformed = time.Now()
		}

		var ctype_transformed int
		switch ctype {
		case "推 ":
			ctype_transformed = 1
		case "噓 ":
			ctype_transformed = -1
		case "→ ":
			ctype_transformed = 0
		default:
			fmt.Printf("Can't transform ctype : %s\n", ctype)
		}

		//fmt.Printf("Find %s(%s)(%s) : %s\n", name, ctime, ctype, content)
		infoList = append(infoList, CommitInfo{
			Type:      ctype_transformed,
			Committer: name,
			Timestamp: ctime_transformed,
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
		PublicUrl:         fromUrl,
		CommitterInfoList: infoList,
		ProcessTime:       time.Now(),
	}

	//fmt.Printf("Prcessed : %s with Committer info list length of : %d", ret.Title, len(ret.CommitterInfoList))

	return ret, nil
}

func createCrawlerClient(fromUrl string) (http.Client, error) {
	var cookies []*http.Cookie
	cookie := &http.Cookie{Name: "over18", Value: "1", Domain: "ptt.cc", Path: "/"}
	cookies = append(cookies, cookie)
	jar, err := cookiejar.New(nil)
	u, err := url.Parse(fromUrl)
	jar.SetCookies(u, cookies)
	client := http.Client{
		Jar: jar,
	}
	return client, err
}

func IterateDocuments(board string, start int, end int, onDocumentPath func(docUrl string)) error {

	if end < start {
		start, end = end, start
	}

	i := start

	for{
		targetUrl := fmt.Sprintf("https://www.ptt.cc/bbs/%s/index%d.html", board, i)
		fmt.Printf("Will parse : %s\n", targetUrl)

		client, _ := createCrawlerClient(targetUrl)
		resp, err := client.Get(targetUrl)
		if err != nil {
			//do something
		}

		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		doc.Find(".r-ent").Each(func(i int, s *goquery.Selection) {
			articleHref, _ := s.Find("a").Attr("href")
			//ret = append(ret, "https://www.ptt.cc/" + articleHref)
			onDocumentPath("https://www.ptt.cc/" + articleHref)
		})

		i++
		if i > end {
			break
		}
	}

	return nil
}

func FetchArticleList(board string, start int, end int) ([]string, error) {

	var ret []string

	err := IterateDocuments(board, start, end, func(docUrl string) {
		ret = append(ret, docUrl)
	})

	return ret, err
}

func ParseRangeDocument(board string, start int, end int) []*PDocRaw {
	var ret []*PDocRaw

	_ = IterateDocuments(board, start, end, func(docUrl string) {
		p, _ := CreateRawDocument(docUrl)
		ret = append(ret, p)
		fmt.Printf("Completed parsed : %s with committer count %d\n", p.Title, len(p.CommitterInfoList))
	})

	fmt.Printf("Completed parsing task : %d", len(ret))
	return ret
}


func Analyze(pdoc *PDocRaw) (*AnalyzedInfo, error) {

	return nil, nil
}
