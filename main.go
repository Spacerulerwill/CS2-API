package main

import (
	"gocasesapi/multiscraper"
	"gocasesapi/util"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Checks if there is any new data to scrape
func main() {
	weaponSkinData := make(map[string]util.Skin)
	caseData := make(map[string]util.Case)
	stickerData := make(map[string]util.Sticker)
	stickerCapsuleData := make(map[string]util.StickerCapsule)
	graffitiData := make(map[string]util.Graffiti)
	musicKitData := make(map[string]util.MusicKit)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Scrape weapon skins
	log.Info().Msg("Scraping weapons skins and knives...")
	weaponSkinLinks, err := util.ReadLines("links/skins.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(weaponSkinLinks, weaponSkinData, 20, multiscraper.ScrapeWeaponLink)

	// Scrape gloves
	log.Info().Msg("Scraping gloves...")
	gloveLinks, err := util.ReadLines("links/gloves.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(gloveLinks, weaponSkinData, 20, multiscraper.ScrapeGloves)

	// Scrape cases
	log.Info().Msg("Scraping cases...")
	caseLinks, err := util.ReadLines("links/cases.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(caseLinks, caseData, 20, multiscraper.ScrapeCase)

	// Scrape stickers
	log.Info().Msg("Scraping stickers...")
	stickerLinks, err := util.ReadLines("links/stickers.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(stickerLinks, stickerData, 20, multiscraper.ScrapeStickerPage)

	// Scrape sticker capsules
	log.Info().Msg("Scraping sticker capsules...")
	stickerCapsuleLinks, err := util.ReadLines("links/capsules.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(stickerCapsuleLinks, stickerCapsuleData, 20, multiscraper.ScrapeStickerCapsule)

	// Scrape single graffitis
	log.Info().Msg("Scraping single grafittis...")
	graffitiPageLinks, err := util.ReadLines("links/graffiti.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(graffitiPageLinks, graffitiData, 20, multiscraper.ScrapeGraffitiPage)

	// Scrape base grade graffiti
	log.Info().Msg("Scraping base grade graffitis...")
	baseGradeGraffitiLinks, err := util.ReadLines("links/base_grade_graffiti.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(baseGradeGraffitiLinks, graffitiData, 20, multiscraper.ScrapeBaseGradeGraffiti)

	// Scrape music kits
	log.Info().Msg("Scraping music kits...")
	musicKitLinks, err := util.ReadLines("links/music_kits.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(musicKitLinks, musicKitData, 20, multiscraper.ScrapeMusicKit)

	// Dump all data to files
	util.WriteJsonToFile("public/api/skins.json", weaponSkinData)
	util.WriteJsonToFile("public/api/cases.json", caseData)
	util.WriteJsonToFile("public/api/stickers.json", stickerData)
	util.WriteJsonToFile("public/api/sticker_capsules.json", stickerCapsuleData)
	util.WriteJsonToFile("public/api/graffiti.json", graffitiData)
	util.WriteJsonToFile("public/api/music_kits.json", musicKitData)
}
