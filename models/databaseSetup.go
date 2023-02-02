package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func SetupModels() {
	var err error

	Dbdriver := os.Getenv("DB_DRIVER")
	dsn := os.Getenv("DB_URL")

	DB, err = gorm.Open(Dbdriver, dsn)
	// DB.LogMode(true)

	if err != nil {
		fmt.Println("Cannot connect to database ", Dbdriver)
		log.Fatal("connection error:", err)
	} else {
		fmt.Println("We are connected to the database ", Dbdriver)
	}

	DB.DropTable("chats")
	DB.DropTable("chat_rooms")
	DB.DropTable("users")

	DB.AutoMigrate(&User{}, &Chat{}, &ChatRoom{})
	DB.Model(&Chat{}).AddForeignKey("sender_name", "users(username)", "CASCADE", "CASCADE")
	DB.Model(&Chat{}).AddForeignKey("receive_name", "users(username)", "CASCADE", "CASCADE")
	DB.Model(&ChatRoom{}).AddForeignKey("sender_name", "users(username)", "CASCADE", "CASCADE")
	DB.Model(&ChatRoom{}).AddForeignKey("receive_name", "users(username)", "CASCADE", "CASCADE")

	SeederUser()
	SeederChatRoom()
	SeederMessages()
}

func SeederUser() {
	//Seeder User
	user1 := User{Username: "anton", Password: "1234", PhotoURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTL_orYH8FKWYS5w45tZdya1Q32e6W0s0ug3g&usqp=CAU"}
	user2 := User{Username: "budi", Password: "1234", PhotoURL: "https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcTCbjYcSIp3PJWZiqcVIXfhyDz4rHklzINnQyEgp4iRzlGZudQ7"}
	user3 := User{Username: "susi", Password: "1234", PhotoURL: "https://pyxis.nymag.com/v1/imgs/176/e83/74be320c0aad12767ee92b95ce29f1a3c4-taylor-swift.1x.rsquare.w1400.jpg"}
	users := []User{user1, user2, user3}

	for _, element := range users {
		err := DB.Create(&element).Error
		if err != nil {
			log.Fatalln("error seeding user")
		}
	}
}

func SeederChatRoom() {
	room1 := ChatRoom{SenderName: "anton", ReceiveName: "budi", LastMessage: "Hai, Anton"}
	room2 := ChatRoom{SenderName: "budi", ReceiveName: "anton", LastMessage: "Hai, Anton"}
	rooms := []ChatRoom{room1, room2}

	for _, element := range rooms {
		err := DB.Create(&element).Error
		if err != nil {
			log.Fatalln("error seeding user")
		}
	}
}

func SeederMessages() {
	message1 := Chat{SenderName: "anton", ReceiveName: "budi", Message: "Hai, Budi"}
	message2 := Chat{SenderName: "budi", ReceiveName: "anton", Message: "Hai, Anton"}

	messages := []Chat{message1, message2}

	for _, element := range messages {
		err := DB.Create(&element).Error
		if err != nil {
			log.Fatalln("error seeding user")
		}
	}
}
