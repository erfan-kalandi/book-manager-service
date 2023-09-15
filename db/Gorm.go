package db

import (
	"errors"
	"fmt"
	"library/config"
	"regexp"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	DBConfig config.Config
	DB       *gorm.DB
}

func InitDB(Config config.Config) (*GormDB, error) {
	//making config and connect to a GormDB
	sqlinfo := fmt.Sprintf("host =%s port=%d user=%s password=%s dbname =%s sslmode=disable",
		Config.Database.Host, Config.Database.Port, Config.Database.Username, Config.Database.Password, Config.Database.DBName)
	DB, err := gorm.Open(postgres.Open(sqlinfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &GormDB{
		DBConfig: Config,
		DB:       DB,
	}, nil
}

func (gdb *GormDB) CreateSchemas() error {
	// Migrate User Accounts
	err := gdb.DB.AutoMigrate(&User{})
	if err != nil {
		return errors.New("can not Migrate User Table")
	}

	//Migrate Book Accounts
	err = gdb.DB.AutoMigrate(&Book{})
	if err != nil {
		return errors.New("can not Migrate Book Table")
	}

	return nil
}

// User
func (gdb *GormDB) CreateNewUser(U *User) error {

	//check username existence
	var count int64
	gdb.DB.Model(&User{}).Where(&User{Username: U.Username}).Count(&count)
	if count > 0 {
		return errors.New("this username is already taken")
	}

	//check email validation
	err := checkmail.ValidateFormat(U.Email)
	if err != nil {
		return errors.New("your Email is not valid")
	}

	// check number validation
	boolean := validatePhone(U.PhoneNumber)
	if !boolean {
		return errors.New("your phone number is not valid")
	}

	// Check email duplication
	gdb.DB.Model(&User{}).Where(&User{Email: U.Email}).Count(&count)
	if count > 0 {
		return errors.New("this email already exists")
	}

	// Check phone number duplication
	gdb.DB.Model(&User{}).Where(&User{PhoneNumber: U.PhoneNumber}).Count(&count)
	if count > 0 {
		return errors.New("this phone number already exists")
	}

	// encrypt the user password
	pass, err := bcrypt.GenerateFromPassword([]byte(U.Password), 10)
	if err != nil {
		return err
	}

	//making new users
	U.Password = string(pass)
	return gdb.DB.Create(U).Error
}

func (gdb *GormDB) GetUserByUsername(Username string) (*User, error) {
	var user User
	err := gdb.DB.Where(&User{Username: Username}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Book
func (gdb *GormDB) AddNewBook(B *Book) error {

	//check Book name existence
	var count int64
	gdb.DB.Model(&Book{}).Where(&Book{Name: B.Name}).Count(&count)
	if count > 0 {
		return errors.New("this book is already exist")
	}

	return gdb.DB.Create(B).Error
}

func (gdb *GormDB) GetAllBooks() (*[]Book, error) {
	var books []Book
	err := gdb.DB.Model(&Book{}).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return &books, nil
}

func (gdb *GormDB) GetBookByID(id int) (*Book, error) {
	var book Book
	err := gdb.DB.Where(&Book{Model: gorm.Model{ID: uint(id)}}).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (gdb *GormDB) DeleteBookByID(id int) error {
	err := gdb.DB.Delete(&Book{Model: gorm.Model{ID: uint(id)}}).Error
	if err != nil {
		return err
	}
	return nil
}

func (gdb *GormDB) UpdateBook(book *Book) error {
	err := gdb.DB.Save(book).Error
	if err != nil {
		return err
	}
	return nil
}

func validatePhone(s string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return re.MatchString(s)
}
