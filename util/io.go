package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

func WriteJsonToFile(filename string, data interface{}) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Failed to marshal json data for file %s", filename))
	}
	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Failed to write file %s", filename))
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
