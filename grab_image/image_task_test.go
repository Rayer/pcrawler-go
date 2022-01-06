package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestExecuteImageStorageTask(t *testing.T) {
	err := ExecuteImageStorageTask(ImageTask{
		BoardName: "Beauty",
		Date:      "1/2/1979",
		Folder:    "testFolder",
		ImageUrl:  "https://i.imgur.com/Ms0bfS9.jpg",
		RootPath:  "/tmp/test",
	})
	fmt.Println(err)
}

func TestExtractImageFromHtmlTask(t *testing.T) {
	u, _ := url.Parse("https://www.ptt.cc/bbs/Beauty/M.1641432468.A.AE6.html")
	fmt.Println(ExtractImageFromHtmlTask(u))
}
