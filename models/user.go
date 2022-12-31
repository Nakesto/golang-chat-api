package models

import (
	"errors"
	"html"
	"math/rand"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nakesto/chat-api/token"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	UID      int    `gorm:"size:255;not null;unique" json:"uid"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	PhotoURL string `gorm:"size:255" json:"photoURL"`
}

func (u *User) SaveUser() (*User, error) {
	err := DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {
	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	u.UID = rand.Intn(10000)

	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (User, string, error) {

	var u User

	err := DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return User{}, "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return User{}, "", err
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return User{}, "", err
	}

	return u, token, nil
}

func GetUserByID(uid uint) (User, error) {

	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}
