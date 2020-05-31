package pcrawler

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type PIndex struct {
	Board      string
	MaxIndex   int
	LastDocUrl *url.URL
	PinnedDocs int
	Articles   []*BriefArticleInfo
}

type BriefArticleInfo struct {
	Title string
	//Can't get precise date because lack of year
	DateString string
	Author     string
	Pinned     bool
	Url        *url.URL
}

func ParseIndexContent(contentIo io.ReadCloser) (*PIndex, error) {

	// Required information :
	// Index count, Last timestamp, Last document url hash

	doc, err := goquery.NewDocumentFromReader(contentIo)

	if err != nil {
		return nil, err
	}

	//Index Count
	maxIndex := 0
	var boardName string
	doc.Find("div#action-bar-container a").Each(func(i int, selection *goquery.Selection) {
		link, exists := selection.Attr("href")
		if !exists {
			return
		}
		r := regexp.MustCompile("/bbs/(.+)/index(\\d+)\\.html")
		if find := r.FindStringSubmatch(link); len(find) > 1 {
			boardName = find[1]
			index, _ := strconv.Atoi(find[2])
			if index+1 > maxIndex {
				maxIndex = index + 1
			}
		}
	})

	pinnedDocs := len(doc.Find(".r-list-sep").NextAll().Nodes)

	var lastDoc *goquery.Selection
	if pinnedDocs > 0 {
		lastDoc = doc.Find(".r-list-sep").Prev()
	} else {
		lastDoc = doc.Find(".r-ent").Last()
	}

	baList := make([]*BriefArticleInfo, 0)
	articleNodes := doc.Find(".r-ent")
	articleCount := len(articleNodes.Nodes)
	articleIndex := 0
	doc.Find(".r-ent").Each(func(i int, s *goquery.Selection) {
		articleHref, _ := s.Find("a").Attr("href")
		articleUrl, _ := url.Parse("https://www.ptt.cc/" + articleHref)
		baList = append(baList, &BriefArticleInfo{
			Title:      s.Find("a").Get(0).FirstChild.Data,
			DateString: strings.Trim(s.Find(".date").Text(), " "),
			Author:     s.Find(".author").Text(),
			Pinned:     articleIndex >= articleCount-pinnedDocs,
			Url:        articleUrl,
		})
		articleIndex++
	})

	lastDocResource, _ := lastDoc.Find(".title a").Attr("href")
	lastDocResource = "https://www.ptt.cc" + lastDocResource
	lastDocUrl, _ := url.Parse(lastDocResource)

	return &PIndex{
		Board:      boardName,
		MaxIndex:   maxIndex,
		LastDocUrl: lastDocUrl,
		PinnedDocs: pinnedDocs,
		Articles:   baList,
	}, nil

}
