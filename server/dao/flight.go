package dao

type Flight struct {
	FlightIdentifier    int32
	SourceLocation      string
	DestinationLocation string
	DepartureTime       int64
	Airfare             float64
	TotalAvailableSeats int32
}
