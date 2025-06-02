package db

import "time"

type RollEvent struct {
	ID    uint      `gorm:"primaryKey"`
	Dtt   time.Time `gorm:"autoCreateTime"`
	Value int       `gorm:"not null;check:value <> 0"`
}
