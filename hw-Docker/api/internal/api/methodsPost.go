package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// rollDice godoc
// @Summary Бросить игральную кость
// @Success 200 {object} responseRoll
// @Tags API V1
// @Router /api/v1/roll_dice [post]
func (handler *apiHandler) rollDice(context *gin.Context) {
	var err error
	var result responseRoll
	result.Result, err = handler.newRollDice()
	if err != nil {
		result.Err = err
		context.JSON(http.StatusInternalServerError, result)
		return
	}
	context.JSON(http.StatusOK, result)
}
