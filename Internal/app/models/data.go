package models

type Data struct {
	Id int     `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}
type Title struct {
	Id string
	X  string
	Y  string
}

const (
	TableName = "measurements"
)
