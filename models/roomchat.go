package models

import (
	"errors"
	"reflect"
	"time"
)

type ChatRoom struct {
	ID          uint       `gorm:"primary_key" json:"-"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"last_message_time"`
	DeletedAt   *time.Time `json:"-" sql:"index"`
	SenderName  string
	ReceiveName string
	LastMessage string `gorm:"size:255" json:"lastmessage"`
	Sender      User   `gorm:"foreignKey:Username;association_foreignkey:SenderName"`
	Receiver    User   `gorm:"foreignKey:Username;association_foreignkey:ReceiveName"`
}

func (chat *ChatRoom) SaveRoom() (*ChatRoom, error) {
	err := DB.Create(&chat).Error
	if err != nil {
		return &ChatRoom{}, err
	}
	return chat, nil
}

func (chat *ChatRoom) BeforeSave() error {
	var cr ChatRoom

	err := DB.Model(chat).Where("(sender_name, receive_name) IN ((?),(?))", chat.SenderName, chat.ReceiveName).Find(&cr).Error

	if err != nil {
		return err
	}

	if reflect.ValueOf(cr).IsNil() {
		return nil
	} else {
		return errors.New("Chat telah dibuat")
	}
}
