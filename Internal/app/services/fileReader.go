package servies

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/Pave1kor/Optics-characteristics/internal/app/models"
)

// Добавить проверку - два ли столбца в файле
func ReadDataFromFile(filePath string) ([]models.Data, models.Title, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, models.Title{}, err
	}
	defer file.Close()

	result := make([]models.Data, 0)
	title := models.Title{}
	scanner := bufio.NewScanner(file)
	// read header
	if scanner.Scan() {
		text := scanner.Text()[0]
		if unicode.IsLetter(rune(text)) {
			header := strings.Fields(scanner.Text())
			title = models.Title{
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
				return nil, models.Title{}, err
			}
		}
		result = append(result, models.Data{
			X: values[0],
			Y: values[1],
		})
	}
	return result, title, scanner.Err()
}
