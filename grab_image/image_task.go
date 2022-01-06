package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Rayer/pcrawler-go"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ImageTask struct {
	BoardName string
	Date      string
	Folder    string
	ImageUrl  string
	RootPath  string
}

type ImageTaskResult struct {
	Task ImageTask
	err  error
}

func ExtractImageFromHtmlTask(url *url.URL) ([]string, error) {
	var ret []string
	client := pcrawler.CreatePttClient(url)
	resp, err := client.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	document.Find("a").Each(func(i int, selection *goquery.Selection) {
		imgUrl, _ := selection.Attr("href")
		if imgUrl != "" {
			if strings.HasSuffix(imgUrl, ".jpg") || strings.HasSuffix(imgUrl, ".jpeg") || strings.HasSuffix(imgUrl, ".png") {
				ret = append(ret, imgUrl)
			}
		}
	})

	return ret, nil
}

func ExecuteImageStorageTask(task ImageTask) ImageTaskResult {

	ret := ImageTaskResult{
		Task: task,
	}

	targetDir := path.Join(task.RootPath, task.BoardName, task.Folder)
	//If doesn't have this folder, create one
	stat, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			ret.err = err
			return ret
		}
	} else if err != nil {
		ret.err = err
		return ret
	}
	if stat != nil && !stat.IsDir() {
		ret.err = fmt.Errorf("path %v is already exist and not a directory", targetDir)
		return ret
	}

	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(task.ImageUrl)
	if err != nil {
		ret.err = err
		return ret
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	filename := filepath.Base(task.ImageUrl)

	file, err := os.Create(path.Join(targetDir, filename))
	if err != nil {
		ret.err = err
		return ret
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, resp.Body)
	ret.err = err

	return ret
}

func ExecuteImageStorageWorker(ctx context.Context, inputChan <-chan ImageTask, outputChan chan<- ImageTaskResult) {

}
