package models

type MessageDTO struct {
	Id       string
	Producer string
	Payload  string
}

type Message struct {
	Id       string
	Producer string
	Payload  string
}
