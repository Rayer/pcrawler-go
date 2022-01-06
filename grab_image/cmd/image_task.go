package cmd

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Rayer/pcrawler-go"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ImageTask struct {
	BoardName string
	Date      string
	Folder    string
	ImageUrl  string
	RootPath  string
}

type RenderResult interface {
	Render() string
}

type ImageTaskResult struct {
	Task ImageTask
	err  error
}

func (i *ImageTaskResult) Render() string {
	return fmt.Sprintf("saving %s to folder %s.... err : %v", i.Task.ImageUrl, i.Task.Folder, i.err)
}

type ImageExtractResult struct {
	TargetUrl *url.URL
	Extracted []string
	err       error
}

func (i *ImageExtractResult) Render() string {
	return fmt.Sprintf("grabbed %v image urls from %s.... err : %v", len(i.Extracted), i.TargetUrl, i.err)
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
			if strings.HasSuffix(imgUrl, ".jpg") || strings.HasSuffix(imgUrl, ".jpeg") ||
				strings.HasSuffix(imgUrl, ".png") || strings.HasSuffix(imgUrl, ".gif") {
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

	client := http.Client{Timeout: 10 * time.Second}
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

func ImageExtractWorker(ctx context.Context, boardName string, rootPath string, inputChan <-chan *pcrawler.BriefArticleInfo, outputChan chan<- ImageTask, renderChan chan<- RenderResult, wg *sync.WaitGroup, wgImgStorage *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-inputChan:
			result := ImageExtractResult{}
			result.TargetUrl = input.Url
			result.Extracted, result.err = ExtractImageFromHtmlTask(input.Url)
			renderChan <- &result
			for _, imageUrl := range result.Extracted {
				wgImgStorage.Add(1)
				outputChan <- ImageTask{
					BoardName: boardName,
					Date:      input.DateString,
					Folder:    input.Title,
					ImageUrl:  imageUrl,
					RootPath:  rootPath,
				}
			}
			wg.Done()
		}
	}
}

func ExecuteImageStorageWorker(ctx context.Context, inputChan <-chan ImageTask, outputChan chan<- RenderResult, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-inputChan:
			result := ExecuteImageStorageTask(input)
			outputChan <- &result
			time.Sleep(time.Duration(rand.Intn(1500)+500) * time.Millisecond)
			wg.Done()
		}
	}
}

func RenderWorker(ctx context.Context, renderInput <-chan RenderResult) {
	for {
		select {
		case <-ctx.Done():
			return
		case render := <-renderInput:
			fmt.Println(render.Render())
		}
	}
}
