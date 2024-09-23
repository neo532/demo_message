package entity

import "time"

type Campaign struct {
	OriginType    int       `json:"-"`
	OriginContent string    `json:"-"`
	MessageType   int       `json:"-"`
	Message       string    `json:"message"`
	TimeSend      time.Time `json:"-"`
}

type Recipient struct {
	Mobile string `json:"mobile"`
	Name   string `json:"name"`
}

type Message struct {
	ID       int64     `json:"id"`
	Status   int       `json:"-"`
	TimeSend time.Time `json:"-"`

	LogID string `json:"-"`

	CampaignID  int64      `json:"-"`
	RecipientID int64      `json:"-"`
	Campaign    *Campaign  `json:"campaign"`
	Recipient   *Recipient `json:"recipient"`
}
