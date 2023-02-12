package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/cyiafn/flight_information_system/server/logs"
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

func GetEnvBool(key string) (bool, bool) {
	envVar := os.Getenv(key)
	if envVar == "" {
		return false, false
	}
	envVar = strings.ToLower(envVar)
	if envVar == "true" {
		return true, true
	}
	return false, true
}
