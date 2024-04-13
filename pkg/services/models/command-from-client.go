package models

type CommandFromClient struct {
	HeartbeatUrl       string `json:"HeartbeatUrl"`
	RPS                string `json:"RPS"`
	TargetRPSTestCheck bool   `json:"TargetRPSCheck"`
}
