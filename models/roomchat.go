package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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

func (room *ChatRoom) SaveRoom() (*ChatRoom, error) {
	err := DB.Create(&room).Error
	if err != nil {
		return &ChatRoom{}, err
	}
	return room, nil
}

func (room *ChatRoom) BeforeSave() error {
	var cr ChatRoom

	if room.SenderName == "" && room.ReceiveName == ""{
		return nil
	}

	fmt.Println(room.SenderName, room.ReceiveName)

	err := DB.Model(ChatRoom{}).Where("(sender_name, receive_name) IN ((?,?))", room.SenderName, room.ReceiveName).Find(&cr).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	
	return errors.New("Chat telah dibuat")
}
