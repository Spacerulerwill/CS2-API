package main

import (
	"fmt"
	"gocasesapi/multiscraper"
	"gocasesapi/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	weaponSkinData     map[string]util.Skin
	caseData           map[string]util.Case
	stickerData        map[string]util.Sticker
	stickerCapsuleData map[string]util.StickerCapsule
	graffitiData       map[string]util.Graffiti
)

func getSkins(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, weaponSkinData)
}

func getCases(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, caseData)
}

func getStickers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, stickerData)
}

func getStickerCapsules(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, stickerCapsuleData)
}

func getGrafitti(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, graffitiData)
}

// Checks if there is any new data to scrape
func updateAPIData() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Scrape weapons
	log.Info().Msg("Scraping skins...")

	weaponSkinLinks, err := util.ReadLines("links/skins.txt")
	if err != nil {
		//log.Fatal(err)
	}

	newWeaponSkinData := make(map[string]util.Skin)
	multiscraper.MultiScrape(weaponSkinLinks, newWeaponSkinData, 20, multiscraper.ScrapeWeaponLink)
	weaponSkinData = newWeaponSkinData

	// Scrape cases
	log.Info().Msg("Scraping cases...")
	caseLinks, err := util.ReadLines("links/cases.txt")
	if err != nil {
		log.Err(err)
	}

	newCaseData := make(map[string]util.Case)
	multiscraper.MultiScrape(caseLinks, newCaseData, 20, multiscraper.ScrapeCase)
	caseData = newCaseData

	// Scrape stickers
	log.Info().Msg("Scraping stickers...")
	stickerLinks, err := util.ReadLines("links/stickers.txt")
	if err != nil {
		log.Err(err)
	}

	newStickerData := make(map[string]util.Sticker)
	multiscraper.MultiScrape(stickerLinks, newStickerData, 20, multiscraper.ScrapeStickerPage)
	stickerData = newStickerData

	// Scrape sticker capsules
	log.Info().Msg("Scraping sticker capsules...")
	stickerCapsuleLinks, err := util.ReadLines("links/capsules.txt")
	if err != nil {
		log.Err(err)
	}

	newStickerCapsuleData := make(map[string]util.StickerCapsule)
	multiscraper.MultiScrape(stickerCapsuleLinks, newStickerCapsuleData, 20, multiscraper.ScrapeStickerCapsule)
	stickerCapsuleData = newStickerCapsuleData

	log.Info().Msg("Scraping single grafittis...")
	newGraffitiData := make(map[string]util.Graffiti)
	graffitiPageLinks, err := util.ReadLines("links/graffiti.txt")
	if err != nil {
		log.Err(err)
	}

	multiscraper.MultiScrape(graffitiPageLinks, newGraffitiData, 20, multiscraper.ScrapeGraffitiPage)

	log.Info().Msg("Scraping base grade graffitis...")
	baseGradeGraffitiLinks, err := util.ReadLines("links/base_grade_graffiti.txt")
	if err != nil {
		log.Err(err)
	}
	fmt.Println(len(baseGradeGraffitiLinks))

	multiscraper.MultiScrape(baseGradeGraffitiLinks, newGraffitiData, 20, multiscraper.ScrapeBaseGradeGraffiti)
	graffitiData = newGraffitiData
}

func main() {
	updateAPIData()

	go func() {
		t := time.NewTicker(time.Hour * 24)
		for {
			<-t.C
			updateAPIData()
		}
	}()

	// Start API
	router := gin.Default()
	router.GET("/skins", getSkins)
	router.GET("/cases", getCases)
	router.GET("/stickers", getStickers)
	router.GET("/sticker-capsules", getStickerCapsules)
	router.GET("/graffiti", getGrafitti)
	router.Run("localhost:8080")
}
