package db

import "time"

// Demo
type Demo struct {
	ID         int64     `gorm:"<-:create;column:id"`
	Name       string    `gorm:"<-;column:name"`
	CreateTime time.Time `gorm:"->;column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"->;column:update_time;autoUpdateTime"`
}
