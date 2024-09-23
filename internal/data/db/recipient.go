package db

// Recipient
type Recipient struct {
	Base

	Mobile string `gorm:"column:mobile"`
	Name   string `gorm:"column:name"`
}
