package handlers

import (
	"context"

	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/dto"
)

func GetFlightInformation(_ context.Context, request any) (any, error) {
	req := request.(*dto.GetFlightInformationRequest)
	res := &dto.GetFlightInformationResponse{}

	foundFlight := false
	for _, flight := range database.GetAllFlights() {
		flight := flight
		if flight.FlightIdentifier == req.FlightIdentifier {
			foundFlight = true
			res.Airfare = flight.Airfare
			res.DepartureTime = flight.DepartureTime
			res.TotalAvailableSeats = flight.TotalAvailableSeats
			break
		}
	}

	if foundFlight {
		return nil, custom_errors.NewNoSuchFlightIdentifierError()
	}

	return res, nil
}
