package main

import (
	"devops-lesta-start-demo/internal/api"
	"devops-lesta-start-demo/internal/config"
	"devops-lesta-start-demo/internal/db"
	"devops-lesta-start-demo/internal/logger"
)

// @BasePath /

// @title DevOps Lesta start demo
// @version 0.0.1
// @description Swagger API for Golang Project DevOps Lesta start demo

// @contact.name API Support
// @contact.email a_guryanov2@lesta.group
func main() {
	logger.InitApiLogger()
	dbClient := db.CreateClient(config.LoadDbConfig(logger.CreateDbLogger()))
	api.InitApi(config.LoadApiConfig(&logger.ApiLog), dbClient)
}
