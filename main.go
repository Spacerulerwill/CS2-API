package main

import (
	"gocasesapi/log"
	"gocasesapi/multiscraper"
	"gocasesapi/util"
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
	patchPackData := make(map[string]util.PatchPack)
	pinData := make(map[string]util.Pin)
	pinCapsuleData := make(map[string]util.PinCapsule)

	// Scrape weapon skins
	log.Info.Println("Scraping weapons skins and knives...")
	weaponSkinLinks, err := util.ReadLines("links/skins.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(weaponSkinLinks, weaponSkinData, 20, multiscraper.ScrapeSkinLink)

	// Scrape gloves
	log.Info.Println("Scraping gloves...")
	gloveLinks, err := util.ReadLines("links/gloves.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(gloveLinks, weaponSkinData, 20, multiscraper.ScrapeGloves)

	// Scrape cases
	log.Info.Println("Scraping cases...")
	caseLinks, err := util.ReadLines("links/cases.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(caseLinks, caseData, 20, multiscraper.ScrapeCase)

	// Scrape collections
	log.Info.Println("Scraping collections...")
	collectionLinks, err := util.ReadLines("links/collections.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(collectionLinks, collectionData, 20, multiscraper.ScrapeCollection)

	log.Info.Println("Scraping souvenir packages...")
	souvenirPackageLinks, err := util.ReadLines("links/souvenir_packages.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(souvenirPackageLinks, souvenirPackageData, 20, multiscraper.ScrapeSouvenirPackagePage)

	// Scrape stickers
	log.Info.Println("Scraping stickers...")
	stickerLinks, err := util.ReadLines("links/stickers.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(stickerLinks, stickerData, 20, multiscraper.ScrapeStickerPage)

	// Scrape sticker capsules
	log.Info.Println("Scraping sticker capsules...")
	stickerCapsuleLinks, err := util.ReadLines("links/capsules.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(stickerCapsuleLinks, stickerCapsuleData, 20, multiscraper.ScrapeStickerCapsule)

	// Scrape single graffitis
	log.Info.Println("Scraping single grafittis...")
	graffitiPageLinks, err := util.ReadLines("links/graffiti.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(graffitiPageLinks, graffitiData, 20, multiscraper.ScrapeGraffitiPage)

	// Scrape base grade graffiti
	log.Info.Println("Scraping base grade graffitis...")
	baseGradeGraffitiLinks, err := util.ReadLines("links/base_grade_graffiti.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(baseGradeGraffitiLinks, graffitiData, 20, multiscraper.ScrapeBaseGradeGraffiti)

	// Scrape music kits
	log.Info.Println("Scraping music kits...")
	musicKitLinks, err := util.ReadLines("links/music_kits.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(musicKitLinks, musicKitData, 20, multiscraper.ScrapeMusicKit)

	// Scrape agents
	log.Info.Println("Scraping agents...")
	agentLinks, err := util.ReadLines("links/agents.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(agentLinks, agentData, 20, multiscraper.ScrapeAgent)

	// Scrape patches
	log.Info.Println("Scraping agent patches...")
	agentPatchesLinks, err := util.ReadLines("links/patches.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(agentPatchesLinks, patchData, 20, multiscraper.ScrapePatch)

	// Scrape patch packs
	log.Info.Println("Scraping patch packs...")
	patchPackLinks, err := util.ReadLines("links/patch_packs.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(patchPackLinks, patchPackData, 20, multiscraper.ScrapePatchPack)

	// Scrape pin pages
	log.Info.Println("Scraping pins...")
	pinPageLinks, err := util.ReadLines("links/pins.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(pinPageLinks, pinData, 20, multiscraper.ScrapePinPage)

	// Scrape pin capsules
	log.Info.Println("Scraping pin capsules...")
	pinCapsuleLinks, err := util.ReadLines("links/pin_capsules.txt")
	if err != nil {
		log.Error.Println(err)
	}
	multiscraper.MultiScrape(pinCapsuleLinks, pinCapsuleData, 20, multiscraper.ScrapePinCapsule)
	multiscraper.WaitForCompletion()
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
	util.WriteJsonToFile("patch_packs.json", patchPackData)
	util.WriteJsonToFile("pins.json", pinData)
	util.WriteJsonToFile("pin_capsules.json", pinCapsuleData)
}
