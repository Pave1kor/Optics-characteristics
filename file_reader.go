package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Energy              float64 `json:"energy"`
	RefractiveIndicator float64 `json:"refractiveIndicator"`
	AbsorptionIndicator float64 `json:"absorptionIndicator"`
}

func ReadDataFromFile(filePath string) ([]Data, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]Data, 0)
	scanner := bufio.NewScanner(file)
	// read header
	if scanner.Scan() {
		header := scanner.Text()
		fmt.Println(header)
	}
	// read data
	for scanner.Scan() {
		dataLine := strings.Fields(scanner.Text())
		values := make([]float64, 3)
		for i := range values {
			if values[i], err = strconv.ParseFloat(dataLine[i], 64); err != nil {
				return nil, err
			}
		}
		result = append(result, Data{
			Energy:              values[0],
			RefractiveIndicator: values[1],
			AbsorptionIndicator: values[2],
		})
	}
	return result, scanner.Err()
}
