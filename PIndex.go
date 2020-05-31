package pcrawler

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"regexp"
	"strconv"
	"time"
)

type IndexOverview struct {
	Board      string
	MaxIndex   int
	LastDocUrl string
	PinnedDocs int
}

type BriefArticleInfo struct {
	Title  string
	Date   time.Time
	Author string
	Pinned bool
}

type PIndexInfo struct {
	IndexOverview
	Articles []*BriefArticleInfo
}

func ParseIndexOverview(contentIo io.ReadCloser) (*IndexOverview, error) {

	// Required information :
	// Index count, Last timestamp, Last document url hash
	ret := IndexOverview{}

	doc, err := goquery.NewDocumentFromReader(contentIo)

	if err != nil {
		return nil, err
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
			if index+1 > ret.MaxIndex {
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

	return &ret, nil
}
