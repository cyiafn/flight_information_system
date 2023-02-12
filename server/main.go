package main

import (
	"github.com/cyiafn/flight_information_system/server/server"
	"github.com/cyiafn/flight_information_system/server/utils"
)

func init() {
	//if err := orm.Init(); err != nil {
	//	panic(err)
	//}
}

func main() {
	defer utils.HandlePanic()
	utils.GracefulShutdown(server.SpinDown)

	server.Boot(routes)
}
