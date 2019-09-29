package PttUtils

import (
	"net/url"
	"testing"
)

func TestCrawler_IndexRawTest(t *testing.T) {
	c, err := NewCrawler("Gossiping")
	if err != nil {
		t.FailNow()
	}
	u, err := url.Parse("https://www.ptt.cc/bbs/Gossiping/index100.html")
	if err != nil {
		t.FailNow()
	}
	content, err := c.getContent(*u)
	if err != nil {
		t.FailNow()
	}
	print(content)
}