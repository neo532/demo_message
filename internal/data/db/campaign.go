package db

import "time"

// Campaign
type Campaign struct {
	Base

	OriginType    int    `gorm:"column:origin_type"`
	OriginContent string `gorm:"column:origin_content"`
	MessageType   int    `gorm:"column:message_type"`
	Message       string `gorm:"column:message"`

	Status   int       `gorm:"column:status"`
	TimeSend time.Time `gorm:"column:time_send"`
}

const (
	CampaignStatusOn  = 1
	CampaignStatusOff = 2
)
