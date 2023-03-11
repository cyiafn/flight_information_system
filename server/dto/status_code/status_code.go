package status_code

import "github.com/cyiafn/flight_information_system/server/custom_errors"

// StatusCodeType are the types possible for errors.
type StatusCodeType uint8

// Do note that under the uber's coding convention, we use iota + 1 as 0 cannot be determined if uninitialised or actually value of 0 in Golang
// The following are error code enums and are self-explanatory.
const (
	// Success returns a successful RPC call.
	Success StatusCodeType = iota + 1

	BusinessLogicGenericError
	MarshallerError

	NoMatchForSourceAndDestination
	NoSuchFlightIdentifier
	InsufficientNumberOfAvailableSeats
)

// GetStatusCode error maps the type of error to the statusCode to return
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
