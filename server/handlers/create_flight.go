package handlers

import (
	"context"

	"github.com/cyiafn/flight_information_system/server/dao"
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/dto"
)

/**
Handlers contains all biz logic of each of the RPC code.
*/

// CreateFlight simply creates a flight and returns to the user the flightIdentifier for the new flight
func CreateFlight(_ context.Context, request any) (any, error) {
	// we cast the request object to the actual object.
	req := request.(*dto.CreateFlightRequest)
	res := &dto.CreateFlightResponse{}

	id := database.GetLargestFlightID() + 1
	database.NewFlight(&dao.Flight{
		FlightIdentifier:    id,
		SourceLocation:      req.SourceLocation,
		DestinationLocation: req.DestinationLocation,
		DepartureTime:       req.DepartureTime,
		Airfare:             req.Airfare,
		TotalAvailableSeats: req.TotalAvailableSeats,
	})

	res.FlightIdentifier = id

	return res, nil
}
