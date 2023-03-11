package handlers

import (
	"context"

	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/dto"
)

// MakeSeatReservation makes a reservation for a flight identifier.
func MakeSeatReservation(_ context.Context, request any) (any, error) {
	req := request.(*dto.MakeSeatReservationRequest)

	foundFlight := false
	for _, flight := range database.GetAllFlights() {
		flight := flight
		if flight.FlightIdentifier != req.FlightIdentifier {
			continue
		}

		foundFlight = true
		if flight.TotalAvailableSeats < req.SeatsToReserve {
			return nil, custom_errors.NewInsufficientNumberOfAvailableSeatsError()
		}
		flight.TotalAvailableSeats -= req.SeatsToReserve

		// handle the callback in the rpc handler in charge of callback
		handleMonitorSeatUpdatesCallback(flight)
	}

	if !foundFlight {
		return nil, custom_errors.NewNoSuchFlightIdentifierError()
	}

	return nil, nil
}
