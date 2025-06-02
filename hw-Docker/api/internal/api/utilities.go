package api

import (
	"math/rand"
)

func (handler *apiHandler) newRollDice() (int, error) {
	value := rand.Intn(6) + 1
	err := handler.dbClient.SaveRollValue(value)
	return value, err
}

func (handler *apiHandler) getStatistic() (Statistic, error) {
	var err error
	result := Statistic{}
	result.AllCountRoll, err = handler.dbClient.GetCount()
	if err != nil {
		return result, err
	}
	result.AllAwgRoll, err = handler.dbClient.GetAwg()
	return result, err
}
