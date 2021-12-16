package entities

import "time"

type Message struct {
	//Id (Auto-Generated)
	Id   string `json:"id"`
	Text string `json:"text"`
	//Date (Auto-Generated)
	Date time.Time `json:"date"`
}
