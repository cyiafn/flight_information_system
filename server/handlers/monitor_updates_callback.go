package handlers

import (
	"context"
	"fmt"
	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/utils/predicates"
	"time"

	"github.com/cyiafn/flight_information_system/server/callback"
	"github.com/cyiafn/flight_information_system/server/dao"
	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/logs"
)

var monitorSeatUpdatesCallbackClient *callback.Client[int32]

func init() {
	monitorSeatUpdatesCallbackClient = callback.NewClient[int32]()
}

func MonitorSeatUpdates(ctx context.Context, request any) (any, error) {
	req := request.(*dto.MonitorSeatUpdatesCallbackRequest)
	exists := predicates.One(database.GetAllFlights(), func(flight *dao.Flight) bool {
		return flight.FlightIdentifier == req.FlightIdentifier
	})
	if !exists {
		return nil, custom_errors.NewNoSuchFlightIdentifierError()
	}

	monitorSeatUpdatesCallbackClient.Subscribe(ctx, req.FlightIdentifier, time.Duration(req.LengthOfMonitorIntervalInSeconds)*time.Second)

	return nil, nil
}

func handleMonitorSeatUpdatesCallback(flight *dao.Flight) {
	res := &dto.MonitorSeatUpdatesCallbackResponse{TotalAvailableSeats: flight.TotalAvailableSeats}
	fmt.Printf("Notifying\n")
	err := monitorSeatUpdatesCallbackClient.Notify(flight.FlightIdentifier, dto.MonitorSeatUpdatesCallbackType, res)
	fmt.Printf("Notified\n")
	if err != nil {
		logs.Warn("failure to deliver callback for 1 or more clients: %v", err)
	}
}
