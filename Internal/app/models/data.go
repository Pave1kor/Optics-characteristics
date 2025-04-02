package models

type DataId struct {
	Date   string
	Number int
}
type TableName struct {
	TableName string
}

type Data struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type Title struct {
	X string `json:"x"`
	Y string `json:"y"`
}
