package models

type ServerLoadTestEvent struct {
	Timestamp         int64 `json:"Timestamp"`
	RPS               int   `json:"RPS"`
	RequestsSucceeded int   `json:"RequestsSucceeded"`
	RequestsFailed    int   `json:"RequestsFailed"`
	NumberOfVUs       int   `json:NumberOfVUs`
}
