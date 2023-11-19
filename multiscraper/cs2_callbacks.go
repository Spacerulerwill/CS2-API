package multiscraper

import (
	"fmt"
	"gocasesapi/log"
	"gocasesapi/util"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var buttonTextIndexMap = map[string]int{
	"Inspect (FN)": 0,
	"Inspect (MW)": 1,
	"Inspect (FT)": 2,
	"Inspect (WW)": 3,
	"Inspect (BS)": 4,
}

type callbackConstraint interface {
	map[string]util.Skin | map[string]util.Sticker | map[string]util.Case |
		map[string]util.StickerCapsule | map[string]util.Graffiti | map[string]util.MusicKit |
		map[string]util.Agent | map[string]util.Patch | map[string]util.Collection |
		map[string]util.SouvenirPackage | map[string]util.PatchPack |
		map[string]util.Pin | map[string]util.PinCapsule
}

// index of condition mapped to the bound of their condition ranges
var (
	skinWearLowerRanges = [5]float32{
		0.45, 0.38, 0.15, 0.07, 0.00,
	}
)

func ScrapeSkinLink(doc *goquery.Document, result map[string]util.Skin, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find(".result-box > h2:nth-child(1)").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	isVanillaKnife := strings.Contains(formattedName, "★ (Vanilla)")
	isDoppler := strings.Contains(formattedName, "Doppler")

	var (
		description, flavorText, minFloatString, maxFloatString, selectedRarity, weaponType string
		minFloat, maxFloat                                                                  float32
		stattrakAvailable, souvenirAvailable                                                bool
		conditionImages, inspectUrls                                                        [5]string
		worstConditionIndex, bestConditionIndex                                             int
	)

	if isVanillaKnife {
		description = ""
		flavorText = ""
		minFloatString = "0.00"
		maxFloatString = "1.00"
		minFloat = 0.00
		maxFloat = 1.00
		selectedRarity = "covert"
		weaponType = "knife"
		stattrakAvailable = true
		souvenirAvailable = false
		worstConditionIndex = 4
		bestConditionIndex = 0

		image := doc.Find(".main-skin-img")
		imageURL, exists := image.Attr("src")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No image URL found for vanilla knife %s", formattedName))
		}

		inspectButton := doc.Find(".inspect-button-skin")
		inspectUrl, exists := inspectButton.Attr("href")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No inspect URL found for vanilla knife %s", formattedName))
		}

		for i := 0; i < 5; i++ {
			conditionImages[i] = imageURL
			inspectUrls[i] = inspectUrl
		}
	} else {

		description = strings.TrimPrefix(doc.Find(".skin-misc-details > p:nth-child(2)").Text(), "Description: ")
		flavorText = doc.Find(".skin-misc-details > p:nth-child(3) > em:nth-child(2) > a:nth-child(1)").Text()

		minFloatString = doc.Find("div.marker-wrapper:nth-child(1) > div:nth-child(1) > div:nth-child(1)").Text()
		{
			minFloat64, err := strconv.ParseFloat(minFloatString, 32)
			if err != nil {
				log.Error.Println(err)
			}

			minFloat = float32(minFloat64)
		}

		maxFloatString = doc.Find("div.marker-wrapper:nth-child(2) > div:nth-child(1) > div:nth-child(1)").Text()
		{
			maxFloat64, err := strconv.ParseFloat(maxFloatString, 32)
			if err != nil {
				log.Error.Println(err)
			}
			maxFloat = float32(maxFloat64)
		}

		// TODO: add check for possibility of failiure
		// Determine best condition index
		for i := 0; i < 5; i++ {
			lowerBound := skinWearLowerRanges[i]
			if minFloat >= lowerBound {
				bestConditionIndex = 4 - i
				break
			}
		}

		for i := 0; i < 5; i++ {
			lowerBound := skinWearLowerRanges[i]
			if maxFloat > lowerBound {
				worstConditionIndex = 4 - i
				break
			}
		}

		skinTypeString := strings.TrimSpace(doc.Find("html body div.container.main-content div.row.text-center div.col-md-10 div.row div.col-md-7.col-widen div.well.result-box.nomargin a.nounderline div p.nomargin").Text())
		stattrakAvailable = doc.Find("div.stattrak").Length() > 0
		souvenirAvailable = doc.Find("div.souvenir").Length() > 0

		rarityFound := false
		for i := 0; i < util.NumSkinRarities; i++ {
			if strings.Contains(skinTypeString, util.SkinRarities[i]) {
				selectedRarity = strings.ToLower(util.SkinRarities[i])
				weaponType = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(skinTypeString, util.SkinRarities[i])))
				rarityFound = true
				break
			}
		}

		if !rarityFound {
			log.Warning.Println(fmt.Sprintf("No rarity found for weapon skin %s", formattedName))
		}

		mainBox := doc.Find("div.well.result-box,nomargin")
		imageButtons := mainBox.Find("a.inspect-button-skin")
		imageButtons.Each(func(i int, button *goquery.Selection) {
			imageURL, exists := button.Attr("data-hoverimg")
			if !exists {
				log.Warning.Println(fmt.Sprintf("No image URL found for weapon skin %s", formattedName))
			}
			inspectUrl, exists := button.Attr("href")
			if !exists {
				log.Warning.Println(fmt.Sprintf("No inspect URL found for weapon skin %s", formattedName))
			}
			index := buttonTextIndexMap[strings.TrimSpace(button.Text())]
			conditionImages[index] = imageURL
			inspectUrls[index] = inspectUrl
		})
	}

	skinData := util.Skin{
		FormattedName:       formattedName,
		Description:         description,
		FlavorText:          flavorText,
		MinFloat:            minFloatString,
		MaxFloat:            maxFloatString,
		WeaponType:          weaponType,
		Rarity:              selectedRarity,
		ConditionImages:     conditionImages,
		WorstConditionIndex: worstConditionIndex,
		BestConditionIndex:  bestConditionIndex,
		InspectUrls:         inspectUrls,
		StattrakAvailable:   stattrakAvailable,
		SouvenirAvailable:   souvenirAvailable,
	}

	if isDoppler {
		skinData.Variations = make(map[string]util.SkinVariation)
		skinVariationWholeBox := doc.Find("#preview-variants > div:nth-child(1)")
		skinBoxes := skinVariationWholeBox.Find("div.no-padding")
		skinBoxes.Each(func(i int, box *goquery.Selection) {
			dopplerFormattedName := box.Find("h3, h4").Text()
			dopplerUnformattedName := util.RemoveNameFormatting(dopplerFormattedName)
			dopplerImage := box.Find("img")
			dopplerImageUrl, exists := dopplerImage.Attr("src")
			if !exists {
				log.Warning.Println(fmt.Sprintf("No image URL found for doppler knife %s %s", formattedName, dopplerFormattedName))
			}
			dopplerInspect := box.Find(".inspect-button-skin")
			dopplerInspectUrl, exists := dopplerInspect.Attr("href")
			if !exists {
				log.Warning.Println(fmt.Sprintf("No inspect URL found for doppler knife %s %s", formattedName, dopplerFormattedName))
			}
			dopplerConditionImages := conditionImages
			dopplerInspectUrls := inspectUrls

			for i := 0; i < 5; i++ {
				if dopplerConditionImages[i] != "" {
					dopplerConditionImages[i] = dopplerImageUrl
					dopplerInspectUrls[i] = dopplerInspectUrl
				}
			}

			skinData.Variations[dopplerUnformattedName] = util.SkinVariation{
				FormattedName:   dopplerFormattedName,
				ConditionImages: dopplerConditionImages,
				InspectUrls:     dopplerInspectUrls,
			}
		})
	}
	mtx.Lock()
	result[unformattedName] = skinData
	mtx.Unlock()
}

func ScrapeCase(doc *goquery.Document, result map[string]util.Case, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find("h1.margin-top-sm").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find(".content-header-img-margin")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No image url for case %s", formattedName))
	}
	skins := make(map[string][]string)
	var rareItemsLink string

	// scrape normal items
	skinBoxes := doc.Find("div.well.result-box.nomargin")
	skinBoxes.Each(func(i int, box *goquery.Selection) {
		rarity := box.Find("div.quality")
		rarityText := rarity.Text()
		if rarityText == "" {
			return
		} else if rarity.HasClass("color-rare-item") {
			rareItemsATag := box.Find("a:nth-child(1)")
			href, exists := rareItemsATag.Attr("href")
			if !exists {
				log.Warning.Println(fmt.Sprintf("No rare items URL found for weapon case %s", formattedName))
			}
			rareItemsLink = href
		} else {
			rarityFound := false
			for i := 0; i < util.NumSkinRarities; i++ {
				if strings.Contains(rarityText, util.SkinRarities[i]) {
					rarityText = strings.ToLower(util.SkinRarities[i])
					rarityFound = true
					break
				}
			}
			skinUnformattedName := util.RemoveNameFormatting(box.Find("h3").Text())

			if !rarityFound {
				log.Warning.Println(fmt.Sprintf("No rarity found for weapon skin %s in weapon case %s", skinUnformattedName, formattedName))
			}

			skins[rarityText] = append(skins[rarityText], skinUnformattedName)
		}
	})

	// scrape rare items
	res := Http2Request(rareItemsLink)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error.Println(err)
	}

	skinBoxes = doc.Find("div.well.result-box.nomargin")
	skinBoxes.Each(func(i int, box *goquery.Selection) {
		rareItemUnformattedName := util.RemoveNameFormatting(box.Find("h3").Text())
		if box.Find("div.quality.color-covert").Length() != 0 {
			skins["rare-item"] = append(skins["rare-item"], rareItemUnformattedName)
		}
	})

	mtx.Lock()
	result[unformattedName] = util.Case{
		FormattedName: formattedName,
		ImageURL:      imageUrl,
		Skins:         skins,
	}
	mtx.Unlock()
}

func ScrapeCollection(doc *goquery.Document, result map[string]util.Collection, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find("div.inline-middle:nth-child(2) > h1:nth-child(1)").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find(".content-header-img-margin")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No image url for collection %s", formattedName))
	}

	// scrape normal items
	skins := make(map[string][]string)
	skinBoxes := doc.Find("div.well.result-box.nomargin")
	skinBoxes.Each(func(i int, box *goquery.Selection) {
		rarity := box.Find("div.quality")
		rarityText := rarity.Text()
		if rarityText == "" {
			return
		} else {
			rarityFound := false
			for i := 0; i < util.NumSkinRarities; i++ {
				if strings.Contains(rarityText, util.SkinRarities[i]) {
					rarityText = strings.ToLower(util.SkinRarities[i])
					rarityFound = true
					break
				}
			}
			skinUnformattedName := util.RemoveNameFormatting(box.Find("h3").Text())

			if !rarityFound {
				log.Warning.Println(fmt.Sprintf("No rarity found for weapon skin %s in collection %s", skinUnformattedName, formattedName))
			}

			skins[rarityText] = append(skins[rarityText], skinUnformattedName)
		}
	})
	mtx.Lock()
	result[unformattedName] = util.Collection{
		FormattedName: formattedName,
		ImageUrl:      imageUrl,
		Skins:         skins,
	}
	mtx.Unlock()
}

func ScrapeSouvenirPackagePage(doc *goquery.Document, result map[string]util.SouvenirPackage, wg *sync.WaitGroup) {
	defer wg.Done()
	boxes := doc.Find("div.well.result-box.nomargin")
	boxes.Each(func(i int, box *goquery.Selection) {
		formattedName := box.Find("h4").Text()
		if formattedName == "" {
			return
		}
		unformattedName := util.RemoveNameFormatting(formattedName)
		image := box.Find("img:nth-child(2)")
		imageUrl, exists := image.Attr("src")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No image url for souvenir package %s", formattedName))
		}

		collection := util.RemoveNameFormatting(box.Find("div:nth-child(1) > div:nth-child(3)").Text())

		mtx.Lock()
		result[unformattedName] = util.SouvenirPackage{
			FormattedName: formattedName,
			ImageURL:      imageUrl,
			Collection:    collection,
		}
		mtx.Unlock()
	})

}

func ScrapeStickerPage(doc *goquery.Document, result map[string]util.Sticker, wg *sync.WaitGroup) {
	defer wg.Done()
	boxes := doc.Find("div.well.result-box.nomargin")
	boxes.Each(func(i int, box *goquery.Selection) {
		formattedName := box.Find("h3").Text()

		if formattedName == "" {
			return
		}

		tournament := box.Find("h4").Text()
		if tournament != "" {
			formattedName += " | " + tournament
		}

		unformattedName := util.RemoveNameFormatting(formattedName)
		image := box.Find("img")
		imageUrl, exists := image.Attr("src")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No image url for sticker %s", formattedName))
		}
		inspectButton := box.Find(".inspect-button-sticker")
		inspectUrl, exists := inspectButton.Attr("href")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No inspect url for sticker %s", formattedName))
		}

		rarityText := box.Find("div.quality").Text()
		rarity := strings.ToLower(strings.TrimSpace(strings.Replace(rarityText, " Sticker", "", 1)))

		mtx.Lock()
		result[unformattedName] = util.Sticker{
			FormattedName: formattedName,
			ImageURL:      imageUrl,
			InspectUrl:    inspectUrl,
			Rarity:        rarity,
		}
		mtx.Unlock()
	})
}

func ScrapeGloves(doc *goquery.Document, result map[string]util.Skin, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find(".result-box > h2:nth-child(1)").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	description := strings.TrimPrefix(doc.Find(".skin-misc-details > p:nth-child(2)").Text(), "Description: ")
	flavorText := doc.Find(".skin-misc-details > p:nth-child(3) > em:nth-child(2)").Text()
	minFloat := doc.Find("div.marker-wrapper:nth-child(1) > div:nth-child(1) > div:nth-child(1)").Text()
	maxFloat := doc.Find("div.marker-wrapper:nth-child(2) > div:nth-child(1) > div:nth-child(1)").Text()
	var conditionImages [5]string
	mainBox := doc.Find("div.well.result-box,nomargin")
	imageButtons := mainBox.Find("a.inspect-button-skin")
	imageButtons.Each(func(i int, button *goquery.Selection) {
		imageURL, exists := button.Attr("data-hoverimg")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No image URL found for weapon skin %s", formattedName))
		}
		index := buttonTextIndexMap[strings.TrimSpace(button.Text())]
		conditionImages[index] = imageURL
	})

	mtx.Lock()
	result[unformattedName] = util.Skin{
		FormattedName:     formattedName,
		Description:       description,
		FlavorText:        flavorText,
		MinFloat:          minFloat,
		MaxFloat:          maxFloat,
		WeaponType:        "gloves",
		Rarity:            "extraordinary",
		ConditionImages:   conditionImages,
		StattrakAvailable: false,
		SouvenirAvailable: false,
	}
	mtx.Unlock()
}

func ScrapeStickerCapsule(doc *goquery.Document, result map[string]util.StickerCapsule, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find("h1.margin-top-sm").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find(".content-header-img-margin")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No image url for sticker capsule %s", formattedName))
	}

	stickers := make(map[string][]string)
	// scrape normal items
	skinBoxes := doc.Find("div.well.result-box.nomargin")
	skinBoxes.Each(func(i int, box *goquery.Selection) {
		rarity := box.Find("div.quality")
		rarityText := rarity.Text()
		if rarityText == "" {
			return
		} else {
			rarityFound := false
			for i := 0; i < util.NumStickerRarities; i++ {
				if strings.Contains(rarityText, util.StickerRarities[i]) {
					rarityText = strings.ToLower(util.StickerRarities[i])
					rarityFound = true
					break
				}
			}
			stickerUnformattedName := util.RemoveNameFormatting(box.Find("h3").Text())

			if !rarityFound {
				log.Warning.Println(fmt.Sprintf("No rarity found for sticker %s in sticker capsule %s", stickerUnformattedName, formattedName))
			}

			stickers[rarityText] = append(stickers[rarityText], stickerUnformattedName)
		}
	})

	mtx.Lock()
	result[unformattedName] = util.StickerCapsule{
		FormattedName: formattedName,
		ImageURL:      imageUrl,
		Stickers:      stickers,
	}
	mtx.Unlock()
}

func ScrapeGraffitiPage(doc *goquery.Document, result map[string]util.Graffiti, wg *sync.WaitGroup) {
	defer wg.Done()
	graffitiBoxes := doc.Find("div.well.result-box.nomargin")
	graffitiBoxes.Each(func(i int, box *goquery.Selection) {
		// formatted and unformatted name
		formattedName := strings.TrimSpace(box.Find("h3").Text())
		if formattedName == "" {
			return
		}
		unformattedName := util.RemoveNameFormatting(formattedName)
		// rarity
		rarityText := box.Find("div.quality").Text()
		if strings.Contains(rarityText, "Base Grade Graffiti") {
			return
		}

		rarityFound := false
		for i := 0; i < util.NumGraffitiRarities; i++ {
			if strings.Contains(rarityText, util.GraffitiRarities[i]) {
				rarityText = strings.ToLower(util.GraffitiRarities[i])
				rarityFound = true
				break
			}
		}

		if !rarityFound {
			log.Warning.Println(fmt.Sprintf("No rarity found for grafitti %s", formattedName))
		}

		// image url
		image := box.Find("img")
		imageUrl, exists := image.Attr("src")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No image url for sticker %s", formattedName))
		}

		inspectButton := box.Find(".inspect-button-graffiti")
		inspectUrl, exists := inspectButton.Attr("href")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No inspect url for sticker %s", formattedName))
		}

		//
		graffitiBox := util.RemoveNameFormatting(box.Find("p.item-resultbox-collection-container-info").Text())

		mtx.Lock()
		result[unformattedName] = util.Graffiti{
			FormattedName: formattedName,
			Rarity:        rarityText,
			ImageURL:      imageUrl,
			InspectUrl:    inspectUrl,
			GraffitiBox:   graffitiBox,
		}
		mtx.Unlock()
	})
}

func ScrapeBaseGradeGraffiti(doc *goquery.Document, result map[string]util.Graffiti, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName :=
		strings.TrimSpace(
			util.RemoveBracketsAndContentsRegex.ReplaceAllString(
				doc.Find(".col-md-8 > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > h2:nth-child(1)").Text(),
				"",
			),
		)

	unformattedName := util.RemoveNameFormatting(formattedName)

	graffitiData := util.Graffiti{
		FormattedName: formattedName,
		Rarity:        "base grade",
	}

	graffitiData.ColorVarations = make(map[string]util.GraffitiColorVariation)

	colorBoxes := doc.Find("div.col-lg-3")
	colorBoxes.Each(func(i int, box *goquery.Selection) {
		color := strings.ToLower(strings.TrimSpace(box.Find("h4").Text()))
		image := box.Find("img")
		imageUrl, exists := image.Attr("src")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No image url for base grade graffiti %s of color %s", formattedName, color))
		}

		inspectButton := box.Find(".inspect-button-graffiti")
		inspectUrl, exists := inspectButton.Attr("href")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No inspect url for base grade graffiti %s of color %s", formattedName, color))
		}
		graffitiData.ColorVarations[color] = util.GraffitiColorVariation{
			ImageUrl:   imageUrl,
			InspectUrl: inspectUrl,
		}
	})

	mtx.Lock()
	result[unformattedName] = graffitiData
	mtx.Unlock()
}

func ScrapeMusicKit(doc *goquery.Document, result map[string]util.MusicKit, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find("div.text-center:nth-child(4) > div:nth-child(1) > div:nth-child(1) > h3:nth-child(1)").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	artist := strings.TrimPrefix(doc.Find("div.text-center:nth-child(4) > div:nth-child(1) > div:nth-child(1) > h4:nth-child(2)").Text(), "By ")
	description := doc.Find(".col-md-8 > div:nth-child(1) > p:nth-child(1)").Text()
	image := doc.Find("img.img-responsive:nth-child(5)")
	stattrakAvailable := doc.Find("div.stattrak").Length() > 0
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No image url for music kit %s", formattedName))
	}
	rarityText := doc.Find("div.quality").Text()
	rarityFound := false
	for i := 0; i < util.NumMusicKitRarities; i++ {
		if strings.Contains(rarityText, util.MusicKitRarities[i]) {
			rarityText = strings.ToLower(util.MusicKitRarities[i])
			rarityFound = true
			break
		}
	}

	if !rarityFound {
		log.Warning.Println(fmt.Sprintf("No rarity found for music kit %s", formattedName))
	}

	musicRows := doc.Find("div.music-file")
	audioUrls := make(map[string]string)
	musicRows.Each(func(i int, box *goquery.Selection) {
		name := box.Find("div:nth-child(1) > p:nth-child(1)").Text()
		audio := box.Find("audio")
		audioSrc, exists := audio.Attr("src")

		if !exists {
			log.Warning.Println(fmt.Sprintf("No audio source for music kit %s %s", formattedName, name))
		}
		audioUrls[name] = "https://csgostash.com" + audioSrc
	})

	var boxesFoundIn []string
	collectionLabels := doc.Find("p.collection-text-label")
	collectionLabels.Each(func(i int, label *goquery.Selection) {
		boxesFoundIn = append(boxesFoundIn, util.RemoveNameFormatting(label.Text()))
	})

	mtx.Lock()
	result[unformattedName] = util.MusicKit{
		FormattedName:     formattedName,
		Artist:            artist,
		Description:       description,
		Rarity:            rarityText,
		ImageURL:          imageUrl,
		StattrakAvailable: stattrakAvailable,
		BoxesFoundIn:      boxesFoundIn,
		AudioURLs:         audioUrls,
	}
	mtx.Unlock()
}

func ScrapeAgent(doc *goquery.Document, result map[string]util.Agent, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find(".col-md-8 > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > h2:nth-child(1)").Text()
	rarityText := doc.Find("div.quality").Text()
	rarityFound := false
	for i := 0; i < util.NumAgentRarities; i++ {
		if strings.Contains(rarityText, util.AgentRarities[i]) {
			rarityText = strings.ToLower(util.AgentRarities[i])
			rarityFound = true
			break
		}
	}

	if !rarityFound {
		log.Warning.Println(fmt.Sprintf("No rarity found for agent %s", formattedName))
	}

	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find("html body div.container.main-content div.row.text-center div.col-md-8.col-widen div.well.result-box.nomargin div.row div.col-md-6 img.img-responsive.center-block.margin-bot-med")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No image url for agent %s", formattedName))
	}

	inspectButton := doc.Find(".inspect-button-pin")
	inspectUrl, exists := inspectButton.Attr("href")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No inspect url for agent %s", formattedName))
	}

	description := strings.TrimPrefix(doc.Find(".skin-misc-details > p:nth-child(1)").Text(), "Description: ")
	flavorText := doc.Find(".col-md-8 > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > p:nth-child(2) > em:nth-child(1)").Text()

	mtx.Lock()
	result[unformattedName] = util.Agent{
		FormattedName: formattedName,
		Rarity:        rarityText,
		ImageUrl:      imageUrl,
		InspectUrl:    inspectUrl,
		Description:   description,
		FlavorText:    flavorText,
	}
	mtx.Unlock()
}

func ScrapePatch(doc *goquery.Document, result map[string]util.Patch, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := strings.TrimPrefix(
		doc.Find(".result-box > div:nth-child(1) > div:nth-child(1) > h2:nth-child(1)").Text(),
		"Patch | ",
	) + " Patch"
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find("img.center-block")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No image url for patch %s", formattedName))
	}

	inspectButton := doc.Find(".inspect-button-pin")
	inspectUrl, exists := inspectButton.Attr("href")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No inspect url for patch %s", formattedName))
	}
	flavorText := doc.Find(".result-box > div:nth-child(2) > div:nth-child(1) > p:nth-child(2) > em:nth-child(1)").Text()
	rarityText := doc.Find("div.quality").Text()
	rarityFound := false
	for i := 0; i < util.NumPatchRarities; i++ {
		if strings.Contains(rarityText, util.PatchRarities[i]) {
			rarityText = strings.ToLower(util.PatchRarities[i])
			rarityFound = true
			break
		}
	}
	if !rarityFound {
		log.Warning.Println(fmt.Sprintf("No rarity found for patch %s", formattedName))
	}

	mtx.Lock()
	result[unformattedName] = util.Patch{
		FormattedName: formattedName,
		Rarity:        rarityText,
		ImageUrl:      imageUrl,
		InspectUrl:    inspectUrl,
		FlavorText:    flavorText,
	}
	mtx.Unlock()
}

func ScrapePatchPack(doc *goquery.Document, result map[string]util.PatchPack, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find(".margin-top-sm").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find(".content-header-img-margin")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No image url for patch pack %s", formattedName))
	}
	patches := make(map[string][]string)

	patchBoxes := doc.Find("div.well.result-box.nomargin")
	patchBoxes.Each(func(i int, box *goquery.Selection) {
		rarity := box.Find("div.quality")
		rarityText := rarity.Text()
		if rarityText == "" {
			return
		} else {
			rarityFound := false
			for i := 0; i < util.NumPatchRarities; i++ {
				if strings.Contains(rarityText, util.PatchRarities[i]) {
					rarityText = strings.ToLower(util.PatchRarities[i])
					rarityFound = true
					break
				}
			}
			patchUnformattedName := util.RemoveNameFormatting(box.Find("h3").Text())

			if !rarityFound {
				log.Warning.Println(fmt.Sprintf("No rarity found for patch %s in patch pack %s", patchUnformattedName, formattedName))
			}

			patches[rarityText] = append(patches[rarityText], patchUnformattedName)
		}
	})
	mtx.Lock()
	result[unformattedName] = util.PatchPack{
		FormattedName: formattedName,
		ImageUrl:      imageUrl,
		Patches:       patches,
	}
	mtx.Unlock()
}

func ScrapePinPage(doc *goquery.Document, result map[string]util.Pin, wg *sync.WaitGroup) {
	defer wg.Done()
	pinBoxes := doc.Find("div.well.result-box.nomargin")
	pinBoxes.Each(func(i int, box *goquery.Selection) {
		formattedName := box.Find("h3").Text()
		if formattedName == "" {
			return
		}

		unformattedName := util.RemoveNameFormatting(formattedName)
		rarityText := box.Find("div.quality").Text()
		rarityFound := false
		for i := 0; i < util.NumCollectibleRarities; i++ {
			if strings.Contains(rarityText, util.CollectibleRarities[i]) {
				rarityText = strings.ToLower(util.CollectibleRarities[i])
				rarityFound = true
				break
			}
		}

		if !rarityFound {
			log.Warning.Println(fmt.Sprintf("No rarity found for pin %s", formattedName))
		}

		image := box.Find("img:nth-child(1)")
		imageUrl, exists := image.Attr("src")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No image url for pin %s", formattedName))
		}

		inspectButton := doc.Find(".inspect-button-pin")
		inspectUrl, exists := inspectButton.Attr("href")
		if !exists {
			log.Warning.Println(fmt.Sprintf("No inspect url for pin %s", formattedName))
		}

		pinCapsule := util.RemoveNameFormatting(box.Find("p.item-resultbox-collection-container-info").Text())

		mtx.Lock()
		result[unformattedName] = util.Pin{
			FormattedName: formattedName,
			Rarity:        rarityText,
			ImageUrl:      imageUrl,
			InspectUrl:    inspectUrl,
			PinCapsule:    pinCapsule,
		}
		mtx.Unlock()
	})
}

func ScrapePinCapsule(doc *goquery.Document, result map[string]util.PinCapsule, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedName := doc.Find(".margin-top-sm").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find("img.content-header-img-margin")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Println(fmt.Sprintf("No image url for pin capsule %s", formattedName))
	}

	pins := make(map[string][]string)

	pinBoxes := doc.Find("div.well.result-box.nomargin")
	pinBoxes.Each(func(i int, box *goquery.Selection) {
		rarity := box.Find("div.quality")
		rarityText := rarity.Text()
		if rarityText == "" {
			return
		} else {
			rarityFound := false
			for i := 0; i < util.NumCollectibleRarities; i++ {
				if strings.Contains(rarityText, util.CollectibleRarities[i]) {
					rarityText = strings.ToLower(util.CollectibleRarities[i])
					rarityFound = true
					break
				}
			}
			pinUnformattedName := util.RemoveNameFormatting(box.Find("h3").Text())

			if !rarityFound {
				log.Warning.Println(fmt.Sprintf("No rarity found for pin %s in pin capsule %s", pinUnformattedName, formattedName))
			}

			pins[rarityText] = append(pins[rarityText], pinUnformattedName)
		}
	})
	mtx.Lock()
	result[unformattedName] = util.PinCapsule{
		FormattedName: formattedName,
		ImageUrl:      imageUrl,
		Pins:          pins,
	}
	mtx.Unlock()
}
