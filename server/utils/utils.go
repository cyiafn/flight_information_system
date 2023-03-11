package utils

import (
	"os"
	"os/signal"
	"runtime/debug"

	json "github.com/bytedance/sonic"
	"github.com/cyiafn/flight_information_system/server/logs"
)

// DumpJSON is an easy way for nested pointers in structures to be printed out in logs
func DumpJSON(a any) string {
	JSON, err := json.MarshalString(a)
	if err != nil {
		return ""
	}
	return JSON
}

// HandlePanic handles panics and prints out the traces on panic
func HandlePanic() {
	if r := recover(); r != nil {
		logs.Error("Recovery Panic: %v, Stack: %s", r, string(debug.Stack()))
	}
}

// GracefulShutdown intercepts SIGINT and SIGKILL and runs cleanup tasks before terminating
func GracefulShutdown(cleanup ...func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		for _, cleanupFunc := range cleanup {
			cleanupFunc()
		}
		os.Exit(0)
	}()

}

// TernaryOperator is a one liner for ternary operations
func TernaryOperator[T any](cond bool, ifTrue, ifFalse T) T {
	if cond {
		return ifTrue
	}
	return ifFalse
}
