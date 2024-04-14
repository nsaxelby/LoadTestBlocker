package models

type ServerEvent struct {
	EventType string `json:"EventType"`
	Data      any    `json:"Data"`
}
