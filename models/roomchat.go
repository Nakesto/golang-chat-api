package models

import (
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

func DeleteRoom(u User, ReceiveName string) (error){
	err := DB.Where("((sender_name, receive_name) IN ((?,?)))", u.Username, ReceiveName).Delete(&ChatRoom{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (room *ChatRoom) SaveRoom() (*ChatRoom, error) {
	err := DB.Create(&room).Error
	if err != nil {
		return &ChatRoom{}, err
	}
	return room, nil
}

func (room *ChatRoom) GetRoom() (*ChatRoom, error) {
	var cr ChatRoom

	err := DB.Unscoped().Model(ChatRoom{}).Where("(sender_name, receive_name) IN ((?,?))", room.SenderName, room.ReceiveName).Find(&cr).Error

	if err != nil {
		return &ChatRoom{}, err
	}

	return &cr, nil
}

func (room *ChatRoom) UpdateDeletedRoom() (error) {

	err := DB.Unscoped().Model(&ChatRoom{}).Where("(sender_name, receive_name) IN ((?,?))", room.SenderName, room.ReceiveName).Update("deleted_at", nil)

	if err != nil {
		return err.Error
	}

	return nil 
}