package handlers

import (
	"context"
	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/utils/predicates"
	"time"

	"github.com/cyiafn/flight_information_system/server/callback"
	"github.com/cyiafn/flight_information_system/server/dao"
	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/logs"
)

// monitorSeatUpdatesCallbackClient is an instance of the callback client for monitoring seats
var monitorSeatUpdatesCallbackClient *callback.Client[int32]

func init() {
	// initialises the client on start
	monitorSeatUpdatesCallbackClient = callback.NewClient[int32]()
}

// MonitorSeatUpdates simply subscribes the client of the RPC call to changes in a particular flight identifier for the time they are provided
func MonitorSeatUpdates(ctx context.Context, request any) (any, error) {
	req := request.(*dto.MonitorSeatUpdatesCallbackRequest)
	// checks if that flight identifier exists
	exists := predicates.One(database.GetAllFlights(), func(flight *dao.Flight) bool {
		return flight.FlightIdentifier == req.FlightIdentifier
	})
	if !exists {
		return nil, custom_errors.NewNoSuchFlightIdentifierError()
	}

	// we subscribe to that flight identifier for changes in seats
	monitorSeatUpdatesCallbackClient.Subscribe(ctx, req.FlightIdentifier, time.Duration(req.LengthOfMonitorIntervalInSeconds)*time.Second)

	return nil, nil
}

// handleMonitorSeatUpdateCallback simply just tells the callback client to notify all subscribers of a flight identifier
func handleMonitorSeatUpdatesCallback(flight *dao.Flight) {
	res := &dto.MonitorSeatUpdatesCallbackResponse{TotalAvailableSeats: flight.TotalAvailableSeats}
	err := monitorSeatUpdatesCallbackClient.Notify(flight.FlightIdentifier, dto.MonitorSeatUpdatesCallbackType, res, nil)
	if err != nil {
		logs.Warn("failure to deliver callback for 1 or more clients: %v", err)
	}
}
