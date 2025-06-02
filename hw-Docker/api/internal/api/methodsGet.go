package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthCheck godoc
// @Summary Проверка доступности сервиса
// @Produce json
// @Success 200 {object} healthCheckState
// @Tags System
// @Router /ping [get]
func (handler *apiHandler) healthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, &health)
}

// statistic godoc
// @Summary Статистика всех бросков кубика
// @Produce json
// @Success 200 {object} responseStat
// @Tags API V1
// @Router /api/v1/roll_statistic [get]
func (handler *apiHandler) statistic(context *gin.Context) {
	var err error
	var result responseStat
	result.Result, err = handler.getStatistic()
	if err != nil {
		result.Err = err
		context.JSON(http.StatusInternalServerError, &result)
		return
	}
	context.JSON(http.StatusOK, &result)
}
