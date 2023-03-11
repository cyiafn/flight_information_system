package handlers

import (
	"context"

	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/dto"
)

// GetFlightIdentifiers simply gets all flight identifiers for a source and destination location
func GetFlightIdentifiers(_ context.Context, request any) (any, error) {
	req := request.(*dto.GetFlightIdentifiersRequest)
	res := &dto.GetFlightIdentifiersResponse{
		FlightIdentifiers: make([]int32, 0),
	}

	for _, flight := range database.GetAllFlights() {
		flight := flight
		if flight.SourceLocation == req.SourceLocation && flight.DestinationLocation == req.DestinationLocation {
			res.FlightIdentifiers = append(res.FlightIdentifiers, flight.FlightIdentifier)
		}
	}

	if len(res.FlightIdentifiers) == 0 {
		return nil, custom_errors.NewNoMatchForSourceAndDestinationError()
	}

	return res, nil
}
