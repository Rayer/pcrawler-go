package test

import (
	"FWFinder"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"regexp"
	"strconv"
	"testing"
)

func TestCrawler_IndexRawTest(t *testing.T) {
	c, err := PttUtils.NewCrawler("Gossiping")
	if err != nil {
		t.FailNow()
	}
	u, err := url.Parse("https://www.ptt.cc/bbs/Gossiping/index100.html")
	if err != nil {
		t.FailNow()
	}
	content, err := c.GetContent(u)
	if err != nil {
		t.FailNow()
	}
	print(content)
}

func TestCrawler_RegexCapture(t *testing.T) {
	link := "/bbs/Gossiping/index12164.html"
	r := regexp.MustCompile("/bbs/.+/index(\\d+)\\.html")
	var index int
	if find := r.FindStringSubmatch(link); len(find) > 1 {
		index, _ = strconv.Atoi(find[1])
		index += 1
	}
	assert.Equal(t, 12165, index)
}

func TestCrawler_GetContent(t *testing.T) {
	c, _ := PttUtils.NewCrawler("Gossiping")
	target, _ := url.Parse("https://www.ptt.cc/bbs/Gossiping/index.html")
	content, _ := c.GetContent(target)
	fmt.Println(content)
}

func TestCrawler_MaxIndex(t *testing.T) {
	c1, _ := PttUtils.NewCrawler("NTUBSE-B-102")
	c2, _ := PttUtils.NewCrawler("Gossiping")

	fmt.Println(c1.GetIndexInitialParameters())
	fmt.Println(c2.GetIndexInitialParameters())

}