package dto

import (
	"github.com/cyiafn/flight_information_system/server/dto/status_code"
	"github.com/cyiafn/flight_information_system/server/logs"
)

type RequestType uint8
type ResponseType uint8

const (
	Ping RequestType = iota + 1
)

var (
	requestToResponseMap = map[RequestType]ResponseType{
		1: 101,
	}
)

func GetResponseType(requestType RequestType) ResponseType {
	res, ok := requestToResponseMap[requestType]
	if !ok {
		logs.Error("Request %v is not mapped to a response", requestType)
	}
	return res
}

type Response struct {
	StatusCode status_code.StatusCodeType
	Data       any
}

func NewRequestDTO(requestType RequestType) any {
	switch requestType {
	case 1:
		return nil
	}
	logs.Error("Request DTO not provided")
	return nil
}
