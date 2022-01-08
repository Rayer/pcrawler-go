package cmd

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Rayer/pcrawler-go"
	"github.com/schollz/progressbar/v3"
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
	BoardName   string
	Date        string
	Folder      string
	ImageUrl    string
	RootPath    string
	NoSubFolder bool
}

func (i *ImageTask) GetDestinationDir() string {
	var ret string
	if noSubFolder {
		ret = path.Join(i.RootPath, i.BoardName)
	} else {
		ret = path.Join(i.RootPath, i.BoardName, i.Folder)
	}

	return ret
}

type ProgressMessage struct {
	AddCrawlerTaskCount     int
	ConsumeCrawlerTaskCount int
	IsCrawlMessageError     bool
	Message                 string
}

type RenderResult interface {
	Progress() ProgressMessage
}

type ImageTaskResult struct {
	Task ImageTask
	err  error
}

func (i *ImageTaskResult) Progress() ProgressMessage {
	return ProgressMessage{
		AddCrawlerTaskCount: 0,
		ConsumeCrawlerTaskCount: func() int {
			if i.err == nil {
				return 1
			}
			return 0
		}(),
		IsCrawlMessageError: i.err != nil,
		Message:             fmt.Sprintf("saving %s into %s...", i.Task.ImageUrl, i.Task.GetDestinationDir()),
	}
}

type ImageExtractResult struct {
	TargetUrl *url.URL
	Extracted []string
	err       error
}

func (i *ImageExtractResult) Progress() ProgressMessage {
	return ProgressMessage{
		AddCrawlerTaskCount:     len(i.Extracted),
		ConsumeCrawlerTaskCount: 0,
		IsCrawlMessageError:     false,
		Message:                 fmt.Sprintf("extracted %v images from %v", len(i.Extracted), i.TargetUrl.String()),
	}
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

func ExecuteImageStorageTask(task ImageTask, noSubFolder bool) ImageTaskResult {

	task.NoSubFolder = noSubFolder
	ret := ImageTaskResult{
		Task: task,
	}

	targetDir := task.GetDestinationDir()

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

func ExecuteImageStorageWorker(ctx context.Context, noSubFolder bool, inputChan <-chan ImageTask, outputChan chan<- RenderResult, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-inputChan:
			result := ExecuteImageStorageTask(input, noSubFolder)
			outputChan <- &result
			time.Sleep(time.Duration(rand.Intn(1500)+500) * time.Millisecond)
			wg.Done()
		}
	}
}

func ShowProgress(ctx context.Context, renderInput <-chan RenderResult) {
	//p := progressbar.Default(100, "working...")
	p := progressbar.NewOptions(1000,
		//progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(false),
		//progressbar.OptionSetWidth(40),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: "[red]-[reset]",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	totalCrawlTask := 0
	consumedCrawlTask := 0
	errCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case render := <-renderInput:
			message := render.Progress()
			totalCrawlTask += message.AddCrawlerTaskCount
			consumedCrawlTask += message.ConsumeCrawlerTaskCount
			if message.IsCrawlMessageError {
				errCount += 1
			}
			showText := message.Message
			if len(showText) > 35 {
				showText = showText[:33] + "..."
			}
			p.ChangeMax(totalCrawlTask)
			_ = p.Set(consumedCrawlTask)
			p.Describe(showText)
		}
	}
}

func ShowLogs(ctx context.Context, renderInput <-chan RenderResult) {
	for {
		select {
		case <-ctx.Done():
			return
		case render := <-renderInput:
			message := render.Progress()
			fmt.Println(message.Message)
		}
	}
}
