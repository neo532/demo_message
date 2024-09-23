package db

import "time"

// Message
type Message struct {
	Base

	LogID string `gorm:"column:log_id"`

	CampaignID  int64     `gorm:"column:campaign_id"`
	RecipientID int64     `gorm:"column:recipient_id"`
	Status      int       `gorm:"column:status"`
	TimeSend    time.Time `gorm:"column:time_send"`

	Campaign  *Campaign  `gorm:"foreignKey:campaign_id;references:id"`
	Recipient *Recipient `gorm:"foreignKey:recipient_id;references:id"`
}

const (
	MessageStatusToSend   = 1
	MessageStatusSended   = 2
	MessageStatusSendFail = 3
)
