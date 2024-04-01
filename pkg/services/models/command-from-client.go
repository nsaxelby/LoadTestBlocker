package models

type CommandFromClient struct {
	HeartbeatUrl string `json:"heartbeaturl"`
	RPS          string `json:"loadrps"`
}
