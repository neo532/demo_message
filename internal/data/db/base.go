package db

import "time"

// Base
type Base struct {
	ID         int64     `gorm:"<-:create;column:id"`
	CreateTime time.Time `gorm:"->;column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"->;column:update_time;autoUpdateTime"`
}
