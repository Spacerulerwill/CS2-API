package main

import (
	"gocasesapi/games/cs2"
	"gocasesapi/log"
	"gocasesapi/multiscraper"
	"gocasesapi/util"
	"os"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func scrapeData[T any](pathToLinks string, outputPath string, callback func(*sync.Mutex, *goquery.Document, map[string]T)) {
	log.Info.Printf("Scraping %s", pathToLinks)
	data := make(map[string]T)
	links, err := util.ReadLines(pathToLinks)
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(links, data, 20, callback)
	util.WriteJsonToFile(outputPath, data)
}

func main() {
	err := os.MkdirAll("output", os.ModePerm)
	if err != nil {
		log.Error.Fatalln(err)
	}

	err = os.MkdirAll("output/cs2", os.ModePerm)
	if err != nil {
		log.Error.Fatalln(err)
	}

	startTime := time.Now()
	scrapeData("links/cs2/skins.txt", "output/cs2/skins.json", cs2.ScrapeSkinLink)
	scrapeData("links/cs2/cases.txt", "output/cs2/cases.json", cs2.ScrapeContainer)
	scrapeData("links/cs2/stickers.txt", "output/cs2/stickers.json", cs2.ScrapeStickerPage)
	scrapeData("links/cs2/sticker_capsules.txt", "output/cs2/sticker_capsules.json", cs2.ScrapeContainer)
	scrapeData("links/cs2/collections.txt", "output/cs2/collections.json", cs2.ScrapeContainer)
	scrapeData("links/cs2/souvenir_packages.txt", "output/cs2/souvenir_packages.json", cs2.ScrapeSouvenirPackagePage)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	log.Info.Printf("Execution time: %s\n", elapsedTime)
}
