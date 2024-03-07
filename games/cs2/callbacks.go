package cs2

import (
	"gocasesapi/log"
	"gocasesapi/util"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var buttonTextIndexMap = map[string]int{
	"Inspect (FN)": 0,
	"Inspect (MW)": 1,
	"Inspect (FT)": 2,
	"Inspect (WW)": 3,
	"Inspect (BS)": 4,
}

var (
	skinWearLowerRanges = [5]float32{
		0.45, 0.38, 0.15, 0.07, 0.00,
	}
)

// Scrape any container from its page
func ScrapeContainer(mtx *sync.Mutex, doc *goquery.Document, result map[string]Container) {
	formattedName := strings.TrimSpace(doc.Find("div.collapsed-top-margin > :nth-child(1)").Text())
	unformattedName := util.RemoveNameFormatting(formattedName)
	image := doc.Find(".content-header-img-margin")
	imageUrl, exists := image.Attr("src")
	if !exists {
		log.Warning.Printf("No image url for %s\n", formattedName)
	}

	// Prefill item ordered map to ensure correct order of rarities
	items := orderedmap.New[string, []string]()
	for _, rarity := range util.QualitiesUnformatted {
		items.Set(rarity, make([]string, 0))
	}

	itemBoxes := doc.Find("div.well.result-box.nomargin")
	itemBoxes.Each(func(i int, box *goquery.Selection) {
		qualityText := box.Find("div.quality").Text()
		// if the box has no rarity text, must be an ad box or something else
		// either way, we can skip it
		if qualityText == "" {
			return
		}

		qualityFound := false
		for _, quality := range util.Qualities {
			if strings.Contains(qualityText, quality) {
				qualityText = strings.ToLower(quality)
				qualityFound = true
				break
			}
		}

		if !qualityFound {
			return
		}

		titleTexts := box.Find("h3, h4")
		var itemUnformattedName string
		titleTexts.Each(func(i int, s *goquery.Selection) {
			itemUnformattedName += s.Text() + " "
		})
		itemUnformattedName = util.RemoveNameFormatting(itemUnformattedName)

		current, _ := items.Get(qualityText)
		items.Set(qualityText, append(current, itemUnformattedName))
	})

	// does the container require a key?
	requiresKey := false
	for _, container := range util.ContainersThatRequireKeys {
		if unformattedName == container {
			requiresKey = true
			break
		}
	}

	// remove empty entries from skin map
	for _, quality := range util.QualitiesUnformatted {
		itemsFromQuality, _ := items.Get(quality)
		if len(itemsFromQuality) == 0 {
			items.Delete(quality)
		}
	}

	mtx.Lock()
	defer mtx.Unlock()
	result[unformattedName] = Container{
		FormattedName: formattedName,
		ImageURL:      imageUrl,
		Items:         items,
		RequiresKey:   requiresKey,
	}
}

// Scraping knives, gloves,
func ScrapeSkinLink(mtx *sync.Mutex, doc *goquery.Document, result map[string]Skin) {
	// Get names
	formattedName := strings.TrimSpace(doc.Find(".result-box > h2:nth-child(1)").Text())
	unformattedName := util.RemoveNameFormatting(formattedName)

	isVanillaKnife := strings.Contains(formattedName, "★ (Vanilla)")
	isDoppler := strings.Contains(formattedName, "Doppler")

	var (
		description, flavorText, minFloatString, maxFloatString, selectedQuality, weaponType string
		minFloat, maxFloat                                                                   float32
		stattrakAvailable, souvenirAvailable                                                 bool
		conditionImages, inspectUrls                                                         [5]string
		worstConditionIndex, bestConditionIndex                                              int
	)

	if isVanillaKnife {
		// Vanilla knives all share the same values for these
		description = ""
		flavorText = ""
		minFloatString = "0.00"
		maxFloatString = "1.00"
		minFloat = 0.00
		maxFloat = 1.00
		selectedQuality = "covert"
		weaponType = "knife"
		stattrakAvailable = true
		souvenirAvailable = false
		worstConditionIndex = 4
		bestConditionIndex = 0

		// Get our single skin image and inspect url
		image := doc.Find(".main-skin-img")
		imageURL, exists := image.Attr("src")
		if !exists {
			log.Warning.Printf("No image URL found for vanilla knife %s\n", formattedName)
		}

		inspectButton := doc.Find(".inspect-button-skin")
		inspectUrl, exists := inspectButton.Attr("href")
		if !exists {
			log.Warning.Printf("No inspect URL found for vanilla knife %s\n", formattedName)
		}

		// Use that single image for all 5 conditions
		for i := 0; i < 5; i++ {
			conditionImages[i] = imageURL
			inspectUrls[i] = inspectUrl
		}
	} else {
		// Determine item rarity and weapon type
		skinTypeString := strings.TrimSpace(doc.Find("div.quality").Text())
		stattrakAvailable = doc.Find("div.stattrak").Length() > 0
		souvenirAvailable = doc.Find("div.souvenir").Length() > 0
		qualityFound := false
		for _, quality := range util.Qualities {
			if strings.Contains(skinTypeString, quality) {
				selectedQuality = strings.ToLower(quality)
				weaponType = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(skinTypeString, quality)))
				qualityFound = true
				break
			}
		}

		if !qualityFound {
			return
		}

		// Description and flavor text come from .misc-details-box
		{
			miscDetailsBox := doc.Find(".skin-misc-details")
			pTags := miscDetailsBox.Find("p")

			// Iterate over each p tag, try to find description
			pTags.EachWithBreak(func(i int, tag *goquery.Selection) bool {
				text := strings.TrimSpace(tag.Text())
				if strings.Contains(text, "Description: ") {
					description = strings.TrimPrefix(text, "Description: ")
					return false
				}
				return true
			})

			// Iterate over each p tag, try to find flavor text
			pTags.EachWithBreak(func(i int, tag *goquery.Selection) bool {
				text := strings.TrimSpace(tag.Text())
				if strings.Contains(text, "Flavor Text: ") {
					flavorText = strings.TrimPrefix(text, "Flavor Text: ")
					return false
				}
				return true
			})
		}

		// Get the min and max floats - keep them as strings for API
		minFloatString = doc.Find("div.marker-wrapper:nth-child(1) > div:nth-child(1) > div:nth-child(1)").Text()
		maxFloatString = doc.Find("div.marker-wrapper:nth-child(2) > div:nth-child(1) > div:nth-child(1)").Text()

		// We still need to convert the min and max float string to actual floats in order
		// to determine the best and worst condition indices
		{
			minFloat64, err := strconv.ParseFloat(minFloatString, 32)
			if err != nil {
				log.Error.Println(err)
			}

			minFloat = float32(minFloat64)
		}

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

		// Get the image urls and inspect urls from the inspect buttons
		mainBox := doc.Find("div.well.result-box.nomargin")
		imageButtons := mainBox.Find("a.inspect-img-hover")
		imageButtons.Each(func(i int, button *goquery.Selection) {
			imageURL, exists := button.Attr("data-hoverimg")
			if !exists {
				log.Warning.Printf("No image URL found for weapon skin %s\n", formattedName)
			}
			inspectUrl, exists := button.Attr("href")
			if !exists {
				log.Warning.Printf("No inspect URL found for weapon skin %s\n", formattedName)
			}
			index := buttonTextIndexMap[strings.TrimSpace(button.Text())]
			conditionImages[index] = imageURL
			inspectUrls[index] = inspectUrl
		})
	}

	// Determine what containers our skin is found in
	containersFoundIn := []string{}
	// only knives and gloves can be found in multiple containers
	// other skins can only be found in one case
	if weaponType == "gloves" || weaponType == "knife" {
		containersList := doc.Find("#knife-cases-collapse > div:nth-child(4)").Find("div")
		containersList.Each(func(i int, s *goquery.Selection) {
			containerFormattedName := strings.TrimSpace(s.Text())
			containerUnformattedName := util.RemoveNameFormatting(containerFormattedName)
			containersFoundIn = append(containersFoundIn, containerUnformattedName)
		})
	} else {
		formattedContainerName := strings.TrimSpace(doc.Find("div.skin-details-collection-container-wrapper:nth-child(1)").Text())
		unformattedContainerName := util.RemoveNameFormatting(formattedContainerName)
		containersFoundIn = append(containersFoundIn, unformattedContainerName)
	}

	// Now detect if the skin has any possible variations
	variations := make(map[string]SkinVariation)
	if isDoppler {
		skinVariationWholeBox := doc.Find("#preview-variants > div:nth-child(1)")
		skinBoxes := skinVariationWholeBox.Find("div.no-padding")
		skinBoxes.Each(func(i int, box *goquery.Selection) {
			dopplerFormattedName := box.Find("h3, h4").Text()
			dopplerUnformattedName := util.RemoveNameFormatting(dopplerFormattedName)
			dopplerImage := box.Find("img")
			dopplerImageUrl, exists := dopplerImage.Attr("src")
			if !exists {
				log.Warning.Printf("No image URL found for doppler knife %s %s\n", formattedName, dopplerFormattedName)
			}
			dopplerInspect := box.Find(".inspect-button-skin")
			dopplerInspectUrl, exists := dopplerInspect.Attr("href")
			if !exists {
				log.Warning.Printf("No inspect URL found for doppler knife %s %s\n", formattedName, dopplerFormattedName)
			}
			dopplerConditionImages := conditionImages
			dopplerInspectUrls := inspectUrls

			for i := 0; i < 5; i++ {
				if dopplerConditionImages[i] != "" {
					dopplerConditionImages[i] = dopplerImageUrl
					dopplerInspectUrls[i] = dopplerInspectUrl
				}
			}

			variations[dopplerUnformattedName] = SkinVariation{
				FormattedName:   dopplerFormattedName,
				ConditionImages: dopplerConditionImages[:],
				InspectUrls:     dopplerInspectUrls[:],
			}
		})
	}

	// Our finished skin data - so beautiful!
	skinData := Skin{
		Item: Item{
			FormattedName:     formattedName,
			Description:       description,
			FlavorText:        flavorText,
			Quality:           selectedQuality,
			InspectURLs:       inspectUrls[:],
			ImageURLs:         conditionImages[:],
			StattrakAvailable: stattrakAvailable,
			SouvenirAvailable: souvenirAvailable,
			ContainersFoundIn: containersFoundIn,
		},
		MinFloat:            minFloatString,
		MaxFloat:            maxFloatString,
		WeaponType:          weaponType,
		WorstConditionIndex: worstConditionIndex,
		BestConditionIndex:  bestConditionIndex,
		Variations:          variations,
	}

	mtx.Lock()
	defer mtx.Unlock()
	result[unformattedName] = skinData
}

// Scrape page full of stickers, don't need to go into the page itself
func ScrapeStickerPage(mtx *sync.Mutex, doc *goquery.Document, result map[string]Sticker) {
	boxes := doc.Find("div.well.result-box.nomargin")
	boxes.Each(func(i int, box *goquery.Selection) {
		formattedName := strings.TrimSpace(box.Find("h3").Text())

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
			log.Warning.Printf("No image url for sticker %s\n", formattedName)
		}
		inspectButton := box.Find(".inspect-button-sticker")
		inspectUrl, exists := inspectButton.Attr("href")
		if !exists {
			log.Warning.Printf("No inspect url for sticker %s\n", formattedName)
		}

		rarityText := box.Find("div.quality").Text()
		rarity := strings.ToLower(strings.TrimSpace(strings.Replace(rarityText, " Sticker", "", 1)))

		container := util.RemoveNameFormatting(strings.TrimSpace(box.Find("p.item-resultbox-collection-container-info").Text()))

		mtx.Lock()
		defer mtx.Unlock()
		result[unformattedName] = Sticker{
			FormattedName:     formattedName,
			Description:       "",
			FlavorText:        "",
			Quality:           rarity,
			InspectURLs:       []string{inspectUrl},
			ImageURLs:         []string{imageUrl},
			StattrakAvailable: false,
			SouvenirAvailable: false,
			ContainersFoundIn: []string{container},
		}
	})
}

// Scrape page full of souvenir packages, don't need to go into the page itself
func ScrapeSouvenirPackagePage(mtx *sync.Mutex, doc *goquery.Document, result map[string]SouvenirPackage) {
	boxes := doc.Find("div.well.result-box.nomargin")
	boxes.Each(func(i int, box *goquery.Selection) {
		formattedName := strings.TrimSpace(box.Find("h4").Text())
		if formattedName == "" {
			return
		}
		unformattedName := util.RemoveNameFormatting(formattedName)
		image := box.Find("img:nth-child(2)")
		imageUrl, exists := image.Attr("src")
		if !exists {
			log.Warning.Printf("No image url for souvenir package %s\n", formattedName)
		}

		collection := util.RemoveNameFormatting(box.Find("div:nth-child(1) > div:nth-child(3)").Text())

		mtx.Lock()
		result[unformattedName] = SouvenirPackage{
			FormattedName: formattedName,
			ImageURL:      imageUrl,
			Collection:    collection,
		}
		mtx.Unlock()
	})
}
