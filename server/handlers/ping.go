package handlers

import (
	"context"

	"github.com/cyiafn/flight_information_system/server/logs"
)

// Ping is a testing function to test connectivity to the server
func Ping(_ context.Context, _ any) (any, error) {
	logs.Info("Received ping, pong...")
	return nil, nil
}
