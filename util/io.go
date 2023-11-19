package util

import (
	"bufio"
	"encoding/json"
	"gocasesapi/log"
	"os"
)

func WriteJsonToFile(filename string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Error.Println(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Error.Println(err)
	}
	_, err = file.Write(jsonData)
	if err != nil {
		log.Error.Println(err)
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
