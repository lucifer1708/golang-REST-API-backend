package models

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func GetUserByID(uid uint) (User, error) {
	var u User
	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := DB.Where("username", username).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func (u *User) SaveUser() (*User, error) {
	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave(*gorm.DB) error {
	hashpaswd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	checkNil(err)
	u.Password = string(hashpaswd)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return nil
}

func VerifyPassword(password, hashpaswd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashpaswd), []byte(password))
}

func (u *User) ValidatePass(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

}

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("databse.db"), &gorm.Config{})
	checkNil(err)
	err = db.AutoMigrate(
		&User{},
	)
	if err != nil {
		return
	}
	DB = db

}

func checkNil(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
