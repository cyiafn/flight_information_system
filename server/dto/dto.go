package dto

import (
	"github.com/cyiafn/flight_information_system/server/dto/status_code"
	"github.com/cyiafn/flight_information_system/server/logs"
)

type RequestType uint8
type ResponseType uint8

const (
	PingRequestType RequestType = iota + 1
	GetFlightIdentifiersRequestType
	GetFlightInformationRequestType
	MakeSeatReservationRequestType
	MonitorSeatUpdatesRequestType
	UpdateFlightPriceRequestType
	CreateFlightRequestType
)

const (
	PingResponseType ResponseType = iota + 101
	GetFlightIdentifiersResponseType
	GetFlightInformationResponseType
	MakeSeatReservationResponseType
	MonitorSeatUpdatesResponseType
	UpdateFlightPriceResponseType
	CreateFlightResponseType
)

const (
	MonitorSeatUpdatesCallbackType = iota + 201
)

var (
	requestToResponseMap = map[RequestType]ResponseType{
		PingRequestType:                 PingResponseType,
		GetFlightIdentifiersRequestType: GetFlightIdentifiersResponseType,
		GetFlightInformationRequestType: GetFlightInformationResponseType,
		MakeSeatReservationRequestType:  MakeSeatReservationResponseType,
		MonitorSeatUpdatesRequestType:   MonitorSeatUpdatesResponseType,
		UpdateFlightPriceRequestType:    UpdateFlightPriceResponseType,
		CreateFlightRequestType:         CreateFlightResponseType,
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
	case PingRequestType:
		return nil
	case GetFlightIdentifiersRequestType:
		return &GetFlightIdentifiersRequest{}
	case GetFlightInformationRequestType:
		return &GetFlightInformationRequest{}
	case MakeSeatReservationRequestType:
		return &MakeSeatReservationRequest{}
	case MonitorSeatUpdatesRequestType:
		return &MonitorSeatUpdatesCallbackRequest{}
	case UpdateFlightPriceRequestType:
		return &UpdateFlightPriceRequest{}
	case CreateFlightRequestType:
		return &CreateFlightRequest{}
	}
	logs.Error("Request DTO not provided")
	return nil
}

type GetFlightIdentifiersRequest struct {
	SourceLocation      string
	DestinationLocation string
}

type GetFlightIdentifiersResponse struct {
	FlightIdentifiers []int32
}

type GetFlightInformationRequest struct {
	FlightIdentifier int32
}

type GetFlightInformationResponse struct {
	DepartureTime       int64
	Airfare             float64
	TotalAvailableSeats int32
}

type MakeSeatReservationRequest struct {
	FlightIdentifier int32
	SeatsToReserve   int32
}

type MonitorSeatUpdatesCallbackRequest struct {
	FlightIdentifier                 int32
	LengthOfMonitorIntervalInSeconds int64
}

type MonitorSeatUpdatesCallbackResponse struct {
	TotalAvailableSeats int32
}

type UpdateFlightPriceRequest struct {
	FlightIdentifier int32
	NewPrice         float64
}

type UpdateFlightPriceResponse struct {
	FlightIdentifier    int32
	SourceLocation      string
	DestinationLocation string
	DepartureTime       int64
	Airfare             float64
	TotalAvailableSeats int32
}

type CreateFlightRequest struct {
	SourceLocation      string
	DestinationLocation string
	DepartureTime       int64
	Airfare             float64
	TotalAvailableSeats int32
}

type CreateFlightResponse struct {
	FlightIdentifier int32
}
