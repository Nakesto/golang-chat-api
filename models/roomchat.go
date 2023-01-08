package models

import "time"

type ChatRoom struct {
	ID          uint       `gorm:"primary_key" json:"-"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"last_message_time"`
	DeletedAt   *time.Time `json:"-" sql:"index"`
	SenderName  string
	ReceiveName string
	LastMessage string `gorm:"size:255" json:"lastmessage"`
	Sender      User   `gorm:"foreignKey:Username;association_foreignkey:SenderName" json:"-"`
	Receiver    User   `gorm:"foreignKey:Username;association_foreignkey:ReceiveName" json:"-"`
}

func (chat *ChatRoom) SaveRoom() (*ChatRoom, error) {
	err := DB.Create(&chat).Error
	if err != nil {
		return &ChatRoom{}, err
	}
	return chat, nil
}
