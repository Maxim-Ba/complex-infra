package models

import "time"

type MessageDTO struct {
	PId       string    `json:"pid" validate:"required"` // id from provider
	Producer  string    `json:"producer" validate:"required"`
	Payload   string    `json:"payload" validate:"required"`
	Group     string    `json:"group" validate:"required"` // чат или комната
	CreatedAt time.Time `json:"сreatedAt"`
}

type WSMessageDTO struct {
	PId       string    `json:"pid" validate:"required"` // id from provider
	Producer  string    `json:"producer" validate:"required"`
	Payload   string    `json:"payload" validate:"required"`
	Group     string    `json:"group" validate:"required"` // чат или комната
	CreatedAt time.Time `json:"сreatedAt"`
	Action    string    `json:"action"`
}
