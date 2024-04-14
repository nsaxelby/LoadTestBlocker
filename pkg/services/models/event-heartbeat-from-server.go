package models

type ServerHeartbeatEvent struct {
	Timestamp int64  `json:"Timestamp"`
	MSLatency int    `json:"MSLatency"`
	Success   bool   `json:"Success"`
	Message   string `json:"Message"`
	Status    string `json:Status`
	Count     int    `json:"Count"`
}
