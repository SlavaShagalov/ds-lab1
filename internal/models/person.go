package models

type Person struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}
