package main

import (
	"github.com/cyiafn/flight_information_system/server/database"
	"github.com/cyiafn/flight_information_system/server/server"
	"github.com/cyiafn/flight_information_system/server/utils"
)

func init() {
	database.PopulateFlights()
}

func main() {
	defer utils.HandlePanic()
	utils.GracefulShutdown(server.SpinDown)

	server.Boot(routes, true)
}
