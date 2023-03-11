package handlers

import (
	"context"

	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/dto"
)

// UpdateFlightPrice updates the flight prices for a particular flight
func UpdateFlightPrice(_ context.Context, request any) (any, error) {
	req := request.(*dto.UpdateFlightPriceRequest)
	res := &dto.UpdateFlightPriceResponse{}

	found := false
	for _, flight := range database.GetAllFlights() {
		flight := flight
		if flight.FlightIdentifier == req.FlightIdentifier {
			found = true
			flight.Airfare = req.NewPrice

			res.FlightIdentifier = req.FlightIdentifier
			res.DepartureTime = flight.DepartureTime
			res.Airfare = flight.Airfare
			res.SourceLocation = flight.SourceLocation
			res.DestinationLocation = flight.DestinationLocation
			res.TotalAvailableSeats = flight.TotalAvailableSeats
			break
		}
	}

	if !found {
		return nil, custom_errors.NewNoSuchFlightIdentifierError()
	}

	return res, nil
}
