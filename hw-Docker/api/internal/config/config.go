package config

import (
	"os"
	"strconv"
	"time"

	"devops-lesta-start-demo/internal/api"
	"devops-lesta-start-demo/internal/db"
	"devops-lesta-start-demo/internal/logger"
)

func getEnvStr(key string, defaultValue string, soft bool) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if soft {
		return defaultValue
	}
	panic(errorEnvNotFound)
}

func getEnvInt(key string, defaultValue int, soft bool) int {
	envValue := getEnvStr(key, strconv.Itoa(defaultValue), soft)
	if envValue == "" && soft {
		return defaultValue
	}
	value, err := strconv.Atoi(envValue)
	if err != nil {
		panic(errorEnvIncorrectFormat)
	}
	return value
}

func LoadApiConfig(logger *logger.ApiGinLogger) api.Conf {
	envPort := getEnvInt("API_PORT", 0, false)
	apiConf := api.Conf{
		Port:   envPort,
		Logger: logger,
	}
	return apiConf
}

func LoadDbConfig(logger *logger.DbGormLogger) db.ConfigDb {
	var timout time.Duration
	var maxConn int
	var maxOpen int
	config := db.ConfigDb{
		DbUrl:  getEnvStr("DB_URL", "", false),
		Logger: logger,
	}
	maxConn = getEnvInt("DB_MAX_CONN", 5, true)
	maxOpen = getEnvInt("DB_MAX_OPEN_CONN", 10, true)
	timoutSecData := getEnvInt("DB_SEC_TIMOUT", 3, true)
	timout = time.Second * time.Duration(timoutSecData)
	config.TimeoutSec = timout
	config.MaxConn = maxConn
	config.MaxOpen = maxOpen
	return config
}
