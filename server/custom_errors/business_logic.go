package custom_errors

import "fmt"

type NoMatchForSourceAndDestinationError struct {
}

func (m *NoMatchForSourceAndDestinationError) Error() string {
	return fmt.Sprintf("No flights match for source and destination locations!")
}

func NewNoMatchForSourceAndDestinationError() error {
	return &NoMatchForSourceAndDestinationError{}
}

type NoSuchFlightIdentifierError struct {
}

func (m *NoSuchFlightIdentifierError) Error() string {
	return fmt.Sprintf("flight identifier provided does not exist")
}

func NewNoSuchFlightIdentifierError() error {
	return &NoSuchFlightIdentifierError{}
}

type InsufficientNumberOfAvailableSeatsError struct {
}

func (m *InsufficientNumberOfAvailableSeatsError) Error() string {
	return fmt.Sprintf("insufficient number of available seats for flight")
}

func NewInsufficientNumberOfAvailableSeatsError() error {
	return &InsufficientNumberOfAvailableSeatsError{}
}
