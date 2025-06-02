package api

import (
	"devops-lesta-start-demo/internal/logger"
)

type db interface {
	SaveRollValue(int) error
	GetCount() (int64, error)
	GetAwg() (float64, error)
}

type apiHandler struct {
	dbClient db
}

type healthCheckState struct {
	Result string `json:"result"`
}

type responseRoll struct {
	Result int   `json:"result"`
	Err    error `json:"error"`
}

type Statistic struct {
	AllCountRoll int64   `json:"all_count_roll"`
	AllAwgRoll   float64 `json:"all_awg_roll"`
}

type responseStat struct {
	Result Statistic `json:"result"`
	Err    error     `json:"error"`
}

type Conf struct {
	Port   int
	Logger *logger.ApiGinLogger
}
