package db

import (
	"context"
	"fmt"

	"devops-lesta-start-demo/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateClient(configDb ConfigDb) *Db {
	dbConn, err := gorm.Open(postgres.Open(configDb.DbUrl), &gorm.Config{Logger: logger.GormLogger})
	if err != nil {
		panic(fmt.Sprintf("failed to create db object, error: %s", err.Error()))
	}
	sqlDb, err := dbConn.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to create sql db object, error: %s", err.Error()))
	}
	sqlDb.SetMaxIdleConns(configDb.MaxConn)
	sqlDb.SetMaxOpenConns(configDb.MaxOpen)
	ctx, cancel := context.WithTimeout(context.Background(), configDb.TimeoutSec)
	defer cancel()
	client := Db{connect: dbConn, context: ctx, log: configDb.Logger}
	migrate(&client)
	return &client
}

func migrate(client *Db) {
	err := client.connect.AutoMigrate(&RollEvent{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database, error: %s", err.Error()))
	}
}

func (client *Db) SaveRollValue(value int) error {
	newData := RollEvent{Value: value}
	result := client.connect.Create(&newData)
	return result.Error
}

func (client *Db) GetCount() (int64, error) {
	var count int64
	result := client.connect.Model(&RollEvent{}).Count(&count)
	if count == 0 {
		return count, errorTableEmpty
	}
	return count, result.Error
}

func (client *Db) GetAwg() (float64, error) {
	var avgValue float64
	result := client.connect.Model(&RollEvent{}).Select("AVG(value)").Scan(&avgValue)
	return avgValue, result.Error
}
