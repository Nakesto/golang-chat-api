package models

import "time"

type Chat struct {
	ID          uint       `gorm:"primary_key" json:"-"`
	CreatedAt   time.Time  `json:"sent_time"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-" sql:"index"`
	Message     string     `gorm:"size:255;not null" json:"message"`
	SenderName  string     `gorm:"size:255;not null" `
	ReceiveName string     `gorm:"size:255;not null" `
	Sender      User       `gorm:"foreignKey:Username;association_foreignkey:SenderName" json:"-"`
	Receiver    User       `gorm:"foreignKey:Username;association_foreignkey:ReceiveName" json:"-"`
}

func (chat *Chat) SaveChat() (*Chat, error) {
	err := DB.Create(&chat).Error
	if err != nil {
		return &Chat{}, err
	}

	return chat, nil
}

func (chat *Chat) BeforeSave() error {
	set1 := []string{chat.SenderName, chat.ReceiveName}
	set2 := []string{chat.ReceiveName, chat.SenderName}

	err := DB.Model(&ChatRoom{}).Where("(sender_name, receive_name) IN ((?),(?))", set1, set2).Update(&ChatRoom{LastMessage: chat.Message}).Error

	if err != nil {
		return err
	}

	return nil
}
