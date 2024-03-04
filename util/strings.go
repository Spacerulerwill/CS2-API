package util

import (
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
var doubleSpaceRegex = regexp.MustCompile(`\s+`)
var RemoveBracketsAndContentsRegex = regexp.MustCompile(`\(.+?\)`)
var substituions = map[string]string{
	"&": "and",
	"ö": "o",
}

var (
	SkinRarities = []string{
		"Consumer Grade",
		"Industrial Grade",
		"Mil-Spec",
		"Restricted",
		"Classified",
		"Covert",
		"Contraband",
	}

	SkinRaritiesUnformatted = make([]string, len(SkinRarities))

	ItemRarities = []string{
		"Base Grade",
		"High Grade",
		"Remarkable",
		"Exotic",
		"Extraordinary",
		"Contraband",
	}

	ItemRaritiesUnformatted = make([]string, len(ItemRarities))

	AgentRarities = []string{
		"Distinguished",
		"Exceptional",
		"Superior",
		"Master",
	}

	SkinConditions = [5]string{
		"Factory New",
		"Minimal Wear",
		"Field Tested",
		"Well-Worn",
		"Battle-Scarred",
	}

	CasesThatDoNotNeedKeys = [2]string{
		"xray p250 package",
		"anubis collection package",
	}
)

func init() {
	for i, rarity := range SkinRarities {
		SkinRaritiesUnformatted[i] = strings.ToLower(rarity)
	}

	for i, rarity := range ItemRarities {
		ItemRaritiesUnformatted[i] = strings.ToLower(rarity)
	}
}
func RemoveNameFormatting(str string) string {
	str = strings.ToLower(str)
	str = nonAlphanumericRegex.ReplaceAllString(str, "")
	str = doubleSpaceRegex.ReplaceAllString(str, " ")
	for k := range substituions {
		str = strings.Replace(str, k, substituions[k], 1)
	}
	return strings.TrimSpace(str)
}
