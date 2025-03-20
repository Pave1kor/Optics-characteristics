package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Energy                float64 `json:"energy"`
	ReflectionCoefficient float64 `json:"reflection_coefficient"`
	AbsorptionCoefficient float64 `json:"absorption_coefficient"`
	Title                 Title   `json:"title"`
}
type Title struct {
	E string
	N string
	K string
}

func readDataFromFile(filePath string) ([]Data, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make([]Data, 0)
	scanner := bufio.NewScanner(file)
	// read header
	if scanner.Scan() {
		// header := strings.Fields(scanner.Text())
		// title = Title{
		// 	E: header[0],
		// 	N: header[1],
		// 	K: header[2],
		// }
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
			Energy:                values[0],
			ReflectionCoefficient: values[1],
			AbsorptionCoefficient: values[2],
		})
	}
	return result, scanner.Err()
}
