package main

import (
	"PttUtils"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"strconv"
)

func main() {
	ap := argparse.NewParser("pdc", "Demo client for FWFinder")

	boardName := ap.String("b", "broadname", &argparse.Options{Required:true, Help:"Board Name"})
	start := ap.String("s", "start", &argparse.Options{Required:true, Help:"Start index"})
	end := ap.String("e", "end", &argparse.Options{Required:true, Help:"End index"})
	err := ap.Parse(os.Args)

	if err != nil {
		fmt.Print(ap.Usage(err))
		os.Exit(1)
	}

	start_i, err := strconv.ParseInt(*start, 10, 0)
	end_i, err := strconv.ParseInt(*end, 10, 0)

	PttUtils.ParseRangeDocument(*boardName, int(start_i), int(end_i))
	//db, err := FWFinder.NewDBObject("node.rayer.idv.tw", "acc", "12qw12qw")
	//if err != nil {
	//	panic("Fail to connect database")
	//}
	//
	//err = FWFinder.IterateDocuments(*boardName, int(start_i), int(end_i), func(docUrl string) {
	//	doc, _ := FWFinder.CreateRawDocument(docUrl)
	//	db.GetDB().
	//
	//})

}
