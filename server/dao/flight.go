package dao

// Flight is the data access object we have for storing information about flights. Generally this will be linked to a db,
// however, for the sake of simplicity, we have stored the data in memory.
type Flight struct {
	FlightIdentifier    int32
	SourceLocation      string
	DestinationLocation string
	DepartureTime       int64
	Airfare             float64
	TotalAvailableSeats int32
}
