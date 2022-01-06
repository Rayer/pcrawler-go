/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/Rayer/pcrawler-go"
	"github.com/spf13/cobra"
	"sync"
)

var startIndex int
var endIndex int
var imageStoragePath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grab_image <board name>",
	Short: "Grab all image in destinated range",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("a board name is required")
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		board := args[0]
		c := pcrawler.NewCrawler(board)
		p, err := c.ParseIndex(0)
		if err != nil {
			return err
		}
		maxIndex := p.MaxIndex
		if startIndex <= 0 {
			startIndex = maxIndex + startIndex
		}
		if endIndex <= 0 {
			endIndex = maxIndex + endIndex
		}

		if startIndex > endIndex {
			startIndex, endIndex = endIndex, startIndex
		}

		articleInfoList := make([]*pcrawler.BriefArticleInfo, 0)

		fmt.Printf("Parsing board %s from %v to %v (most recent : %v)...\n", board, startIndex, endIndex, maxIndex)
		for i := startIndex; i <= endIndex; i++ {
			indexInfo, err := c.ParseIndex(i)
			if err != nil {
				fmt.Println("err : ", err)
				continue
			}
			articleInfoList = append(articleInfoList, indexInfo.Articles...)
		}

		renderCtx, renderCancel := context.WithCancel(context.Background())

		imgExtractCtx, imgExtractCancel := context.WithCancel(context.Background())
		imgStorageCtx, imgStorageCancel := context.WithCancel(context.Background())
		wgImgExtract := &sync.WaitGroup{}
		wgImgStorage := &sync.WaitGroup{}
		wgImgExtract.Add(len(articleInfoList))
		imgExtractInput := make(chan *pcrawler.BriefArticleInfo)
		renderChannel := make(chan RenderResult, 200)
		imgTaskChannel := make(chan ImageTask, 200)
		for i := 0; i < 5; i++ {
			go ImageExtractWorker(imgExtractCtx, board, imageStoragePath, imgExtractInput, imgTaskChannel, renderChannel, wgImgExtract, wgImgStorage)
		}
		for i := 0; i < 5; i++ {
			go ExecuteImageStorageWorker(imgStorageCtx, imgTaskChannel, renderChannel, wgImgStorage)
		}
		go RenderWorker(renderCtx, renderChannel)
		for _, a := range articleInfoList {
			imgExtractInput <- a
		}

		wgImgExtract.Wait()
		wgImgStorage.Wait()
		imgExtractCancel()
		imgStorageCancel()
		renderCancel()

		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVarP(&startIndex, "start", "s", -10, "")
	rootCmd.Flags().IntVarP(&endIndex, "end", "e", 0, "")
	rootCmd.Flags().StringVarP(&imageStoragePath, "path", "p", "", "")

	_ = rootCmd.MarkFlagRequired("path")

}
