package util

import (
	"bufio"
	"os"
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

var SkinRarities = [7]string{
	"Consumer Grade",
	"Industrial Grade",
	"Mil-Spec",
	"Restricted",
	"Classified",
	"Covert",
	"Contraband",
}

var StickerRarities = [5]string{
	"High Grade",
	"Remarkable",
	"Exotic",
	"Extraordinary",
	"Contraband",
}

var GraffitiRarities = [4]string{
	"Base Grade",
	"High Grade",
	"Remarkable",
	"Exotic",
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

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
