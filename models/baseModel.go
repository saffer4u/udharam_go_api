package models

type Base struct {
	IsSuccess     bool   `json:"isSuccess"`
	Code          int    `json:"code"`
	StatusMessage string `json:"statusMessage"`
	Data          any    `json:"data"`
}

var Response *Base

func CreateResponse(IsSuccess bool, Code int, StatusMessage string, Data any) {
	Response = &Base{
		IsSuccess:     IsSuccess,
		Code:          Code,
		StatusMessage: StatusMessage,
		Data:          Data,
	}

}
