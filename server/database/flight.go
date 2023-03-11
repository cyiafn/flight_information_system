package database

import "github.com/cyiafn/flight_information_system/server/dao"

/*
Note: Query is a linear scan, it is not efficient, but this is not the focus of this project
*/

// flights is our in-memory "db" to store flights
var flights []*dao.Flight

// this emulates the auto-incrementing PK function of some databases
var largestFlightID int32

// GetAllFlights Getter function to return all flights
func GetAllFlights() []*dao.Flight {
	return flights
}

// GetLargestFlightID gets the largest flight, as there is no flight delete functionality, the following is sufficient
func GetLargestFlightID() int32 {
	return largestFlightID
}

// NewFlight emulates an insert with an auto-incrementing PK.
func NewFlight(flight *dao.Flight) {
	flights = append(flights, flight)
	largestFlightID += 1
}

// PopulateFlights simply populates hardcoded data for flights.
func PopulateFlights() {
	flights = append(flights, &dao.Flight{
		FlightIdentifier:    1,
		SourceLocation:      "Singapore",
		DestinationLocation: "San Francisco",
		DepartureTime:       1701388800,
		Airfare:             2050.6,
		TotalAvailableSeats: 99,
	})
	flights = append(flights, &dao.Flight{
		FlightIdentifier:    2,
		SourceLocation:      "Singapore",
		DestinationLocation: "San Francisco",
		DepartureTime:       1701388900,
		Airfare:             3239.20,
		TotalAvailableSeats: 54,
	})
	flights = append(flights, &dao.Flight{
		FlightIdentifier:    3,
		SourceLocation:      "Singapore",
		DestinationLocation: "Kular Lumpur",
		DepartureTime:       1701287800,
		Airfare:             99.9,
		TotalAvailableSeats: 22,
	})
	flights = append(flights, &dao.Flight{
		FlightIdentifier:    4,
		SourceLocation:      "Singapore",
		DestinationLocation: "Bali",
		DepartureTime:       1701176800,
		Airfare:             325.1,
		TotalAvailableSeats: 2,
	})
	flights = append(flights, &dao.Flight{
		FlightIdentifier:    5,
		SourceLocation:      "Tokyo",
		DestinationLocation: "Seoul",
		DepartureTime:       1701065800,
		Airfare:             892.2,
		TotalAvailableSeats: 1,
	})
	flights = append(flights, &dao.Flight{
		FlightIdentifier:    6,
		SourceLocation:      "Tokyo",
		DestinationLocation: "Shanghai",
		DepartureTime:       1701054800,
		Airfare:             239.2,
		TotalAvailableSeats: 2,
	})
	largestFlightID = 6
}
