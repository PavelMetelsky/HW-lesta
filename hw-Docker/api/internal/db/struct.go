package db

import (
	"context"
	"time"

	"devops-lesta-start-demo/internal/logger"

	"gorm.io/gorm"
)

type Db struct {
	connect *gorm.DB
	context context.Context
	log     *logger.DbGormLogger
}

type ConfigDb struct {
	DbUrl      string
	MaxConn    int
	MaxOpen    int
	TimeoutSec time.Duration
	Logger     *logger.DbGormLogger
}
