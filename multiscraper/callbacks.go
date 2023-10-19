package multiscraper

import (
	"fmt"
	"gocasesapi/util"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
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
		map[string]util.Agent
}

func ScrapeWeaponLink(doc *goquery.Document, result map[string]util.Skin) {
	formattedName := doc.Find(".result-box > h2:nth-child(1)").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	isVanillaKnife := strings.Contains(formattedName, "★ (Vanilla)")
	isDoppler := strings.Contains(formattedName, "Doppler")

	var (
		description, flavorText, minFloat, maxFloat, selectedRarity, weaponType string
		stattrakAvailable, souvenirAvailable                                    bool
		conditionImages                                                         [5]string
	)

	if isVanillaKnife {
		description = ""
		flavorText = ""
		minFloat = "0.00"
		maxFloat = "1.00"
		selectedRarity = "covert"
		weaponType = "knife"
		stattrakAvailable = true
		souvenirAvailable = false

		image := doc.Find(".main-skin-img")
		imageURL, exists := image.Attr("src")

		if !exists {
			log.Warn().Msg(fmt.Sprintf("No image URL found for vanilla knife %s", formattedName))
		}

		for i := 0; i < 5; i++ {
			conditionImages[i] = imageURL
		}
	} else {
		description = strings.TrimPrefix(doc.Find(".skin-misc-details > p:nth-child(2)").Text(), "Description: ")
		flavorText = doc.Find(".skin-misc-details > p:nth-child(3) > em:nth-child(2) > a:nth-child(1)").Text()
		minFloat = doc.Find("div.marker-wrapper:nth-child(1) > div:nth-child(1) > div:nth-child(1)").Text()
		maxFloat = doc.Find("div.marker-wrapper:nth-child(2) > div:nth-child(1) > div:nth-child(1)").Text()
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
			log.Warn().Msg(fmt.Sprintf("No rarity found for weapon skin %s", formattedName))
		}

		mainBox := doc.Find("div.well.result-box,nomargin")
		imageButtons := mainBox.Find("a.inspect-button-skin")
		imageButtons.Each(func(i int, button *goquery.Selection) {
			imageURL, exists := button.Attr("data-hoverimg")
			if !exists {
				log.Warn().Msg(fmt.Sprintf("No image URL found for weapon skin %s", formattedName))
			}
			index := buttonTextIndexMap[strings.TrimSpace(button.Text())]
			conditionImages[index] = imageURL
		})
	}

	skinData := util.Skin{
		FormattedName:     formattedName,
		Description:       description,
		FlavorText:        flavorText,
		MinFloat:          minFloat,
		MaxFloat:          maxFloat,
		WeaponType:        weaponType,
		Rarity:            selectedRarity,
		ConditionImages:   conditionImages,
		StattrakAvailable: stattrakAvailable,
		SouvenirAvailable: souvenirAvailable,
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
				log.Warn().Msg(fmt.Sprintf("No image URL found for doppler knife %s %s", formattedName, dopplerFormattedName))
			}
			dopplerConditionImages := conditionImages
			for i := 0; i < 5; i++ {
				if dopplerConditionImages[i] != "" {
					dopplerConditionImages[i] = dopplerImageUrl
				}
			}

			skinData.Variations[dopplerUnformattedName] = util.SkinVariation{
				FormattedName:   dopplerFormattedName,
				ConditionImages: dopplerConditionImages,
			}
		})
	}
	mtx.Lock()
	result[unformattedName] = skinData
	mtx.Unlock()
}

func ScrapeCase(doc *goquery.Document, result map[string]util.Case) {
	formattedName := doc.Find("h1.margin-top-sm").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find(".content-header-img-margin")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warn().Msg(fmt.Sprintf("No image url for case %s", formattedName))
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
				log.Warn().Msg(fmt.Sprintf("No rare items URL found for weapon case %s", formattedName))
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
				log.Warn().Msg(fmt.Sprintf("No rarity found for weapon skin %s in weapon case %s", skinUnformattedName, formattedName))
			}

			skins[rarityText] = append(skins[rarityText], skinUnformattedName)
		}
	})

	// scrape rare items
	res := Http2Request(rareItemsLink)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Err(err)
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

func ScrapeStickerPage(doc *goquery.Document, result map[string]util.Sticker) {
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
			log.Warn().Msg(fmt.Sprintf("No image url for sticker %s", formattedName))
		}
		rarityText := box.Find("div.quality").Text()
		rarity := strings.ToLower(strings.TrimSpace(strings.Replace(rarityText, " Sticker", "", 1)))

		mtx.Lock()
		result[unformattedName] = util.Sticker{
			FormattedName: formattedName,
			ImageURL:      imageUrl,
			Rarity:        rarity,
		}
		mtx.Unlock()
	})
}

func ScrapeGloves(doc *goquery.Document, result map[string]util.Skin) {
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
			log.Warn().Msg(fmt.Sprintf("No image URL found for weapon skin %s", formattedName))
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

func ScrapeStickerCapsule(doc *goquery.Document, result map[string]util.StickerCapsule) {
	formattedName := doc.Find("h1.margin-top-sm").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find(".content-header-img-margin")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warn().Msg(fmt.Sprintf("No image url for sticker capsule %s", formattedName))
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
				log.Warn().Msg(fmt.Sprintf("No rarity found for sticker %s in sticker capsule %s", stickerUnformattedName, formattedName))
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

func ScrapeGraffitiPage(doc *goquery.Document, result map[string]util.Graffiti) {
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
			log.Warn().Msg(fmt.Sprintf("No rarity found for grafitti %s", formattedName))
		}

		// image url
		image := box.Find("img")
		imageUrl, exists := image.Attr("src")
		if !exists {
			log.Warn().Msg(fmt.Sprintf("No image url for sticker %s", formattedName))
		}

		//
		graffitiBox := util.RemoveNameFormatting(box.Find("p.item-resultbox-collection-container-info").Text())

		mtx.Lock()
		result[unformattedName] = util.Graffiti{
			FormattedName: formattedName,
			Rarity:        rarityText,
			ImageURL:      imageUrl,
			GraffitiBox:   graffitiBox,
		}
		mtx.Unlock()
	})
}

func ScrapeBaseGradeGraffiti(doc *goquery.Document, result map[string]util.Graffiti) {
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

	graffitiData.ColorVarations = make(map[string]string)

	colorBoxes := doc.Find("div.col-lg-3")
	colorBoxes.Each(func(i int, box *goquery.Selection) {
		color := strings.ToLower(strings.TrimSpace(box.Find("h4").Text()))
		image := box.Find("img")
		imageUrl, exists := image.Attr("src")
		if !exists {
			log.Warn().Msg(fmt.Sprintf("No image url for base grade graffiti %s of color %s", formattedName, color))
		}
		graffitiData.ColorVarations[color] = imageUrl
	})

	mtx.Lock()
	result[unformattedName] = graffitiData
	mtx.Unlock()
}

func ScrapeMusicKit(doc *goquery.Document, result map[string]util.MusicKit) {
	formattedName := doc.Find("div.text-center:nth-child(4) > div:nth-child(1) > div:nth-child(1) > h3:nth-child(1)").Text()
	unformattedName := util.RemoveNameFormatting(formattedName)
	artist := strings.TrimPrefix(doc.Find("div.text-center:nth-child(4) > div:nth-child(1) > div:nth-child(1) > h4:nth-child(2)").Text(), "By ")
	description := doc.Find(".col-md-8 > div:nth-child(1) > p:nth-child(1)").Text()
	image := doc.Find("img.img-responsive:nth-child(5)")
	stattrakAvailable := doc.Find("div.stattrak").Length() > 0
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warn().Msg(fmt.Sprintf("No image url for music kit %s", formattedName))
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
		log.Warn().Msg(fmt.Sprintf("No rarity found for music kit %s", formattedName))
	}

	musicRows := doc.Find("div.music-file")
	audioUrls := make(map[string]string)
	musicRows.Each(func(i int, box *goquery.Selection) {
		name := box.Find("div:nth-child(1) > p:nth-child(1)").Text()
		audio := box.Find("audio")
		audioSrc, exists := audio.Attr("src")

		if !exists {
			log.Warn().Msg(fmt.Sprintf("No audio source for music kit %s %s", formattedName, name))
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

func ScrapeAgent(doc *goquery.Document, result map[string]util.Agent) {
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
		log.Warn().Msg(fmt.Sprintf("No rarity found for agent %s", formattedName))
	}

	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find("html body div.container.main-content div.row.text-center div.col-md-8.col-widen div.well.result-box.nomargin div.row div.col-md-6 img.img-responsive.center-block.margin-bot-med")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warn().Msg(fmt.Sprintf("No image url for music kit %s", formattedName))
	}

	description := strings.TrimPrefix(doc.Find(".skin-misc-details > p:nth-child(1)").Text(), "Description: ")
	flavorText := doc.Find(".col-md-8 > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > p:nth-child(2) > em:nth-child(1)").Text()

	mtx.Lock()
	result[unformattedName] = util.Agent{
		FormattedName: formattedName,
		Rarity:        rarityText,
		ImageUrl:      imageUrl,
		Description:   description,
		FlavorText:    flavorText,
	}
	mtx.Unlock()
}
