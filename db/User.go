package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `gorm:"varchar(50),unique"`
	Firstname   string `gorm:"varchar(50)"`
	Lastname    string `gorm:"varchar(50)"`
	PhoneNumber string `gorm:"varchar(30),unique"`
	Email       string `gorm:"varchar(100),unique"`
	Gender      string `gorm:"varchar(50)"`
	Password    string `gorm:"varchar(64)"`
}
