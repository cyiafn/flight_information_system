package utils

import (
	"os"
	"os/signal"
	"runtime/debug"

	json "github.com/bytedance/sonic"
	"github.com/cyiafn/flight_information_system/server/logs"
)

func DumpJSON(a any) string {
	JSON, err := json.MarshalString(a)
	if err != nil {
		return ""
	}
	return JSON
}

func Ptr[T any](a T) *T {
	return &a
}

func HandlePanic() {
	if r := recover(); r != nil {
		logs.Error("Recovery Panic: %v, Stack: %s", r, string(debug.Stack()))
	}
}

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

func TernaryOperator[T any](cond bool, ifTrue, ifFalse T) T {
	if cond {
		return ifTrue
	}
	return ifFalse
}
