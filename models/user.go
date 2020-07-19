package models

import (
	"time"

	"github.com/jinzhu/gorm"
	drivers "github.com/zbcheng/filestore/drivers/mysql"
)

var db *gorm.DB

func init() {
	db = drivers.DBConn()
}

type User struct {
	gorm.Model
	Username   string    `gorm:"username" json:"username"`
	Password   string    `gorm:"password" json:"password"`
	Email      string    `gorm:"email" json:"email"`
	Phone      string    `gorm:"phone" json:"phone"`
	SignupAt   time.Time `gorm:"signup_at datetime" json:"signup_at"`
	LastActive time.Time `gorm:"last_active datetime" json:"last_active"`
	Status     int       `gorm:"status" json:"status"`
	Token      string    `gorm:"token" json:"token"`
}

func (user *User) UpdateLastActive() {
	lastActive := time.Now()
	db.Update(&user).UpdateColumn("last_active", lastActive)
}

func (user *User) UpdateToken(token string) {
	user.Token = token
	db.Save(&user)
}
