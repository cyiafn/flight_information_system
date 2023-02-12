package main

import (
	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/handlers"
)

var routes = map[dto.RequestType]func(request any) (any, error){
	1: handlers.Ping,
}
