package util

import (
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
var doubleSpaceRegex = regexp.MustCompile(`\s+`)
var RemoveBracketsAndContentsRegex = regexp.MustCompile("\\(.+?\\)")
var substituions = map[string]string{
	"&": "and",
	"ö": "o",
}

const (
	NumSkinRarities     int = 7
	NumStickerRarities  int = 5
	NumGraffitiRarities int = 4
	NumMusicKitRarities int = 1
)

var SkinRarities = [NumSkinRarities]string{
	"Consumer Grade",
	"Industrial Grade",
	"Mil-Spec",
	"Restricted",
	"Classified",
	"Covert",
	"Contraband",
}

var StickerRarities = [NumStickerRarities]string{
	"High Grade",
	"Remarkable",
	"Exotic",
	"Extraordinary",
	"Contraband",
}

var GraffitiRarities = [NumGraffitiRarities]string{
	"Base Grade",
	"High Grade",
	"Remarkable",
	"Exotic",
}

var MusicKitRarities = [NumMusicKitRarities]string{
	"High Grade",
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
