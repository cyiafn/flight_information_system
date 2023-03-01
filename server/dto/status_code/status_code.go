package status_code

import "github.com/cyiafn/flight_information_system/server/custom_errors"

type StatusCodeType uint8

const (
	Success StatusCodeType = iota + 1

	BusinessLogicGenericError
	MarshallerError

	NoMatchForSourceAndDestination
	NoSuchFlightIdentifier
	InsufficientNumberOfAvailableSeats
)

func GetStatusCode(err error) StatusCodeType {
	if err == nil {
		return Success
	}
	switch err.(type) {
	case *custom_errors.MarshallerError:
		return MarshallerError
	case *custom_errors.NoMatchForSourceAndDestinationError:
		return NoMatchForSourceAndDestination
	case *custom_errors.NoSuchFlightIdentifierError:
		return NoSuchFlightIdentifier
	case *custom_errors.InsufficientNumberOfAvailableSeatsError:
		return InsufficientNumberOfAvailableSeats
	default:
		return BusinessLogicGenericError
	}
}
