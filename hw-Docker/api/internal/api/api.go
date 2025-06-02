package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func InitApi(apiConf Conf, dbClient db) {
	routes := gin.New()
	routes.Use(*apiConf.Logger.Filters)
	routes.Use(gin.Recovery())
	handler := &apiHandler{
		dbClient: dbClient,
	}
	metric := ginprometheus.NewPrometheus("gin")
	metric.Use(routes)
	registerRoutes(handler, routes)
	err := routes.Run(":" + strconv.Itoa(apiConf.Port))
	panic(err)
}
