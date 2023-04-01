package model

type Response struct {
	Status string
	Error  string
	Code   int
	Data   interface{}
}
