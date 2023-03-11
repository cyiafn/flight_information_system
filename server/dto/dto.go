package dto

import (
	"github.com/cyiafn/flight_information_system/server/dto/status_code"
	"github.com/cyiafn/flight_information_system/server/logs"
)

/**
This file contains everything to do with the data transfer objects.
Generally, we will use an IDL such as Protobuf with the compiler such as protoc to generate the DTO objects with the mappings
from request to response. If we use a RPC frameworks such as gRPC, the handler objects and mapping will be generated for us too.

However, we have done it manually in this file for simplicity.
*/

// RequestType is the type of request object in the payload
type RequestType uint8

// ResponseType is the type of response object in the payload
type ResponseType uint8

// Each of this request types correspond with a request for an RPC call. 1 - 100 are requests
const (
	PingRequestType RequestType = iota + 1
	GetFlightIdentifiersRequestType
	GetFlightInformationRequestType
	MakeSeatReservationRequestType
	MonitorSeatUpdatesRequestType
	UpdateFlightPriceRequestType
	CreateFlightRequestType
)

// Each of these response types correspond with a response for an RPC call. 101 - 200 are responses
const (
	PingResponseType ResponseType = iota + 101
	GetFlightIdentifiersResponseType
	GetFlightInformationResponseType
	MakeSeatReservationResponseType
	MonitorSeatUpdatesResponseType
	UpdateFlightPriceResponseType
	CreateFlightResponseType
)

// MonitorSeatUpdatesCallbackType Each of these callback types correspond with a callback for a subscription. 201 - 300 are callback messsages
const (
	MonitorSeatUpdatesCallbackType = iota + 201
)

var (
	// requestToResponseMap simply maps the request to the relevant response types
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

// GetResponseType simply maps the request type to the appropriate response type
func GetResponseType(requestType RequestType) ResponseType {
	res, ok := requestToResponseMap[requestType]
	if !ok {
		logs.Error("Request %v is not mapped to a response", requestType)
	}
	return res
}

// Response is a generic wrapper around any response object. Data contains the actual payload of the output of the RPC call
// while StatusCode contains the status of the RPC call. Note that Data will be nil in the event that StatusCode != 1
type Response struct {
	StatusCode status_code.StatusCodeType
	Data       any
}

// NewRequestDTO generates a new
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

/*
The following are request and response data transfer objects defined for each RPC call.

Some RPC calls may have no response body and only return a statusCode
*/

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
