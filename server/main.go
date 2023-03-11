package main

import (
	"context"
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/handlers"
	"github.com/cyiafn/flight_information_system/server/server"
	"github.com/cyiafn/flight_information_system/server/utils"
)

// startup initialisation
func init() {
	database.PopulateFlights()
}

// entry point
func main() {
	// handles panics
	defer utils.HandlePanic()
	// spinsdown server upon terminating application
	utils.GracefulShutdown(server.SpinDown)
	// boots up the server
	server.Boot(routes, true)
}

// routes are the  routes from request to handlers
var routes = map[dto.RequestType]func(ctx context.Context, request any) (any, error){
	dto.PingRequestType:                 handlers.Ping,
	dto.GetFlightIdentifiersRequestType: handlers.GetFlightIdentifiers,
	dto.GetFlightInformationRequestType: handlers.GetFlightInformation,
	dto.MakeSeatReservationRequestType:  handlers.MakeSeatReservation,
	dto.MonitorSeatUpdatesRequestType:   handlers.MonitorSeatUpdates,
	dto.UpdateFlightPriceRequestType:    handlers.UpdateFlightPrice,
	dto.CreateFlightRequestType:         handlers.CreateFlight,
}
