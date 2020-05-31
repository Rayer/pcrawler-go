package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"pcrawler"
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

	pcrawler.ParseRangeDocument(*boardName, int(start_i), int(end_i))

}
