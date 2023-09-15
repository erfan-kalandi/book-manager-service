package db

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name              string `gorm:"varchar(50),unique"`
	Category          string `gorm:"varchar(50)"`
	Volume            int
	PublishedAt       time.Time `gorm:"TIMESTAMP"`
	Summary           string    `gorm:"varchar(250)"`
	PublisherName     string    `gorm:"varchar(50)"`
	Owner             string    `gorm:"varchar(50)"`
	TableOfContents   []string  `gorm:"json"`
	AuthorFirstName   string    `gorm:"varchar(50)"`
	AuthorLastName    string    `gorm:"varchar(50)"`
	AuthorNationality string    `gorm:"varchar(50)"`
	AuthorBirthday    time.Time `gorm:"TIMESTAMP"`
}
