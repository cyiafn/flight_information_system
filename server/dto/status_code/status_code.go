package status_code

import "github.com/cyiafn/flight_information_system/server/custom_errors"

type StatusCodeType uint8

const (
	Success StatusCodeType = iota + 1
	BusinessLogicError
	DatabaseError
	NetworkError
	MarshallerError
)

func GetStatusCode(err error) StatusCodeType {
	if err == nil {
		return Success
	}
	switch err.(type) {
	case *custom_errors.MarshallerError:
		return MarshallerError
	}
	return Success
}
