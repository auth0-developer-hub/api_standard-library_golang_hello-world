package models

type Metadata struct {
	Api 	string `json:"api"`
	Branch 	string `json:"branch"`
}

type ApiResponse struct {
	Metadata Metadata `json:"metadata"`
	Text string `json:"text"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}
