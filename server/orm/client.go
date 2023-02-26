package orm

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	defaultHost = "localhost"
	hostEnvKey  = "DB_HOST"

	defaultUser = "postgres"
	userEnvKey  = "DB_USER"

	defaultPassword = "admin"
	passwordEnvKey  = "DB_PASSWORD"

	defaultDBName = "binance-futures-leaderboard"
	dbNameEnvKey  = "DB_NAME"

	defaultPort = "5432"
	portEnvKey  = "DB_PORT"

	defaultSSLMode = false
	sslModeEnvKey  = "DB_SSL_ENABLED"

	defaultTimeZone = "Asia/Singapore"
	timeZoneEnvKey  = "DB_TZ"
)

var client *gorm.DB

func Init() error {
	var err error
	client, err = gorm.Open(postgres.Open(makeDSN()), &gorm.Config{})
	if err != nil {
		logs.Error("unable to connect to db, err: %v", err)
		return err
	}
	return nil
}

func makeDSN() string {
	host := defaultHost
	envVarHost, ok := utils.GetEnvStr(hostEnvKey)
	if ok {
		host = envVarHost
	}

	user := defaultUser
	envVarUser, ok := utils.GetEnvStr(userEnvKey)
	if ok {
		user = envVarUser
	}

	password := defaultPassword
	envVarPassword, ok := utils.GetEnvStr(passwordEnvKey)
	if ok {
		password = envVarPassword
	}

	dbName := defaultDBName
	envVarDBName, ok := utils.GetEnvStr(dbNameEnvKey)
	if ok {
		dbName = envVarDBName
	}

	port := defaultPort
	envVarPort, ok := utils.GetEnvStr(portEnvKey)
	if ok {
		port = envVarPort
	}

	sslMode := getSSLMode(defaultSSLMode)
	envVarSSLMode, ok := utils.GetEnvBool(sslModeEnvKey)
	if ok {
		sslMode = getSSLMode(envVarSSLMode)
	}

	timeZone := defaultTimeZone
	envVarTimeZone, ok := utils.GetEnvStr(timeZoneEnvKey)
	if ok {
		timeZone = envVarTimeZone
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbName, port, sslMode, timeZone)
}

func getSSLMode(enabled bool) string {
	if enabled {
		return "enable"
	}
	return "disable"
}

func GetClient() *gorm.DB {
	return client
}

func CloseConnOnExit() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigs
		db, _ := client.DB()
		db.Close()
		os.Exit(0)
	}()
}
