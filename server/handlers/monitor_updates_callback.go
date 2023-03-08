package handlers

import (
	"context"
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
	monitorSeatUpdatesCallbackClient.Subscribe(ctx, req.FlightIdentifier, time.Duration(req.LengthOfMonitorIntervalInSeconds)*time.Second)

	return nil, nil
}

func handleMonitorSeatUpdatesCallback(flight *dao.Flight) {
	res := &dto.MonitorSeatUpdatesCallbackResponse{TotalAvailableSeats: flight.TotalAvailableSeats}
	err := monitorSeatUpdatesCallbackClient.Notify(flight.FlightIdentifier, dto.MonitorSeatUpdatesResponseType, res)
	if err != nil {
		logs.Warn("failure to deliver callback for 1 or more clients: %v", err)
	}
}
