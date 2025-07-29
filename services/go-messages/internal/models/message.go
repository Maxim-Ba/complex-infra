package models

type MessageDTO struct {
	Id       string `json:"id"`
	Producer string `json:"producer"`
	Payload  string `json:"payload"`
	Group    string `json:"group"` // чат или комнат
}

type Message struct {
	Id       string
	Producer string
	Payload  string
	Group    string
}
