package models

type DataId struct {
	measurement_date   string
	measurement_number int
}

type Data struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type Title struct {
	X string `json:"x"`
	Y string `json:"y"`
}
