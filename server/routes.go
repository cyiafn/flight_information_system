package main

import (
	"context"

	"github.com/cyiafn/flight_information_system/server/dto"
	"github.com/cyiafn/flight_information_system/server/handlers"
)

var routes = map[dto.RequestType]func(ctx context.Context, request any) (any, error){
	1: handlers.Ping,
}
