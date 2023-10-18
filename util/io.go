package util

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

func WriteJsonToFile(filename string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Err(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Err(err)
	}
	_, err = file.Write(jsonData)
	if err != nil {
		log.Err(err)
	}
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
