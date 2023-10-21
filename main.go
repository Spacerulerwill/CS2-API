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
	collectionData := make(map[string]util.Collection)
	souvenirPackageData := make(map[string]util.SouvenirPackage)
	stickerData := make(map[string]util.Sticker)
	stickerCapsuleData := make(map[string]util.StickerCapsule)
	graffitiData := make(map[string]util.Graffiti)
	musicKitData := make(map[string]util.MusicKit)
	agentData := make(map[string]util.Agent)
	patchData := make(map[string]util.Patch)

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

	// Scrape collections
	log.Info().Msg("Scraping collections...")
	collectionLinks, err := util.ReadLines("links/collections.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(collectionLinks, collectionData, 20, multiscraper.ScrapeCollection)

	log.Info().Msg("Scraping souvenir packages...")
	souvenirPackageLinks, err := util.ReadLines("links/souvenir_packages.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(souvenirPackageLinks, souvenirPackageData, 20, multiscraper.ScrapeSouvenirPackagePage)

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

	// Scrape agents
	log.Info().Msg("Scraping agents...")
	agentLinks, err := util.ReadLines("links/agents.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(agentLinks, agentData, 20, multiscraper.ScrapeAgent)

	log.Info().Msg("Scraping agent patches...")
	agentPatchesLinks, err := util.ReadLines("links/patches.txt")
	if err != nil {
		log.Err(err)
	}
	multiscraper.MultiScrape(agentPatchesLinks, patchData, 20, multiscraper.ScrapePatch)

	// Dump all data to files
	util.WriteJsonToFile("skins.json", weaponSkinData)
	util.WriteJsonToFile("cases.json", caseData)
	util.WriteJsonToFile("collections.json", collectionData)
	util.WriteJsonToFile("souvenir_packages.json", souvenirPackageData)
	util.WriteJsonToFile("stickers.json", stickerData)
	util.WriteJsonToFile("sticker_capsules.json", stickerCapsuleData)
	util.WriteJsonToFile("graffiti.json", graffitiData)
	util.WriteJsonToFile("music_kits.json", musicKitData)
	util.WriteJsonToFile("agents.json", agentData)
	util.WriteJsonToFile("patches.json", patchData)
}
