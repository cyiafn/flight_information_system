package utils

import (
	"github.com/cyiafn/flight_information_system/server/logs"
	"os"
	"strconv"
)

func GetEnvStr(key string) (string, bool) {
	envVar := os.Getenv(key)
	if envVar == "" {
		return "", false
	}
	return envVar, true
}

func GetEnvInt(key string) (int, bool) {
	envVar := os.Getenv(key)
	if envVar == "" {
		return 0, false
	}
	intEnvVar, err := strconv.Atoi(envVar)
	if err != nil {
		logs.Warn("unable to convert envVar to int, val: %s, err: %v", envVar, err)
		return 0, false
	}
	return intEnvVar, false
}
