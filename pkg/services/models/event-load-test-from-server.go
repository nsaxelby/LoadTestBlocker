package models

type ServerLoadTestEvent struct {
	Timestamp int64 `json:"Timestamp"`
	RPS       int   `json:"RPS"`
	Count     int   `json:"Count"`
	VU        int   `json:VU`
}
