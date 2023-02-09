package models

type Entry struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Response struct {
	Result any `json:"result"`
}
