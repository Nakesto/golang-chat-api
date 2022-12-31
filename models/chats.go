package models

import "github.com/jinzhu/gorm"

type Chat struct {
	gorm.Model
	Message     string `gorm:"size:255;not null" json:"message"`
	SenderName  string
	ReceiveName string
	Sender      User `gorm:"foreignKey:SenderName;references:Username"`
	Receiver    User `gorm:"foreignKey:ReceiveName;references:Username"`
}

func (chat *Chat) SaveChat() (*Chat, error) {
	err := DB.Create(&chat).Error
	if err != nil {
		return &Chat{}, err
	}

	return chat, nil
}
