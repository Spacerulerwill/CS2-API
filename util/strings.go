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
	Qualities = []string{
		"Consumer Grade",
		"Industrial Grade",
		"Mil-Spec",
		"Restricted",
		"Classified",
		"Covert",
		"Contraband",
		"Base Grade",
		"High Grade",
		"Remarkable",
		"Exotic",
		"Extraordinary",
		"Distinguished",
		"Exceptional",
		"Superior",
		"Master",
	}

	QualitiesUnformatted = make([]string, len(Qualities))

	SkinConditions = [5]string{
		"Factory New",
		"Minimal Wear",
		"Field Tested",
		"Well-Worn",
		"Battle-Scarred",
	}

	ContainersThatRequireKeys = []string{
		"winter offensive weapon case",
		"spectrum 2 case",
		"spectrum case",
		"snakebite case",
		"shattered web case",
		"shadow case",
		"revolver case",
		"prisma 2 case",
		"prisma case",
		"operation wildfire case",
		"operation vanguard weapon case",
		"operation riptide case",
		"operation pheonix weapon case",
		"operation hydra case",
		"operation broken fang case",
		"operation breakout weapon case",
		"operation bravo case",
		"huntsman weapon case",
		"horizon case",
		"glove case",
		"gamma 2 case",
		"gamma case",
		"falchion case",
		"esports 2014 summer case",
		"esports 2013 winter case",
		"esports 2013 case",
		"danger zone case",
		"cs20 case",
		"csgo weapon case 3",
		"csgo weapon case 2",
		"csgo weapon case",
		"clutch case",
		"chroma 3 case",
		"chroma 2 case",
		"chroma case",
		"revolution case",
		"recoil case",
		"fracture case",
		"dreams and nightmares case",
		"kilowatt case",
		"sticker capsule 2",
		"sticker capsule",
	}
)

func init() {
	for i, rarity := range Qualities {
		QualitiesUnformatted[i] = strings.ToLower(rarity)
	}

}
func RemoveNameFormatting(str string) string {
	str = strings.ToLower(str)
	for k := range substituions {
		str = strings.Replace(str, k, substituions[k], 1)
	}
	str = nonAlphanumericRegex.ReplaceAllString(str, "")
	str = doubleSpaceRegex.ReplaceAllString(str, " ")
	return strings.TrimSpace(str)
}
