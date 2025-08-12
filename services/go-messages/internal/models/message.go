package models

import "time"

type MessageDTO struct {
	PId        string    `json:"pid"` // id from provider
	Producer  string    `json:"producer"`
	Payload   string    `json:"payload"`
	Group     string    `json:"group"` // чат или комната
	CreatedAt time.Time `json:"сreatedAt"`
}

type Message struct {
	Id       string
	Producer string
	Payload  string
	Group    string
}


type RequestMessages struct {
	GroupiD string 
	Offset int32 // смещение от конца массива
	Count int32 //сколько сообщений запрашивается

}
