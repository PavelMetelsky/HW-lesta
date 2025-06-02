package api

import (
	_ "devops-lesta-start-demo/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerRoutes(handler *apiHandler, routes *gin.Engine) {
	routes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.GET("/ping", handler.healthCheck)
	apiV1 := routes.Group("/api/v1/")
	apiV1.POST("roll_dice", handler.rollDice)
	apiV1.GET("roll_statistic", handler.statistic)
}
