package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Data struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type Title struct {
	X string `json:"x"`
	Y string `json:"y"`
}

// Добавить проверку - два ли столбца в файле
func readDataFromFile(filePath string) ([]Data, Title, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, Title{}, err
	}
	defer file.Close()

	result := make([]Data, 0)
	title := Title{}
	scanner := bufio.NewScanner(file)
	// read header
	if scanner.Scan() {
		text := scanner.Text()[0]
		if unicode.IsLetter(rune(text)) {
			header := strings.Fields(scanner.Text())
			title = Title{
				X: header[0],
				Y: header[1],
			}
		}
	}
	// read data
	for scanner.Scan() {
		dataLine := strings.Fields(scanner.Text())
		values := make([]float64, 2)
		for i := range values {
			if values[i], err = strconv.ParseFloat(dataLine[i], 64); err != nil {
				return nil, Title{}, err
			}
		}
		result = append(result, Data{
			X: values[0],
			Y: values[1],
		})
	}
	return result, title, scanner.Err()
}
