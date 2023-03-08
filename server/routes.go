package main

import (
	"context"

	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/handlers"
)

var routes = map[dto.RequestType]func(ctx context.Context, request any) (any, error){
	dto.PingRequestType:                 handlers.Ping,
	dto.GetFlightIdentifiersRequestType: handlers.GetFlightIdentifiers,
	dto.GetFlightInformationRequestType: handlers.GetFlightInformation,
	dto.MakeSeatReservationRequestType:  handlers.MakeSeatReservation,
	dto.MonitorSeatUpdatesRequestType:   handlers.MonitorSeatUpdates,
	dto.UpdateFlightPriceRequestType:    handlers.UpdateFlightPrice,
	dto.CreateFlightRequestType:         handlers.CreateFlight,
}
