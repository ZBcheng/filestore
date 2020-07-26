package models

import (
	"github.com/jinzhu/gorm"
	drivers "github.com/zbcheng/filestore/app/drivers/mysql"
)

var db *gorm.DB

func init() {
	db = drivers.DBConn()
}

type User struct {
	ID         int    `gorm:"id" json:"id"`
	Username   string `gorm:"username" json:"username"` // 用户名
	Password   string `gorm:"password" json:"password"` // 密码
	Email      string `gorm:"email" json:"email"`       // 邮箱
	Phone      string `gorm:"phone" json:"phone"`       // 手机
	Avatar     string `gorm:"avatar" json:"avatar"`     // 头像
	Status     int    `gorm:"status" json:"status"`
	Token      string `gorm:"token" json:"token"`
	SignupAt   string `gorm:"signup_at" json:"signup_at"`
	LastActive string `gorm:"last_active" json:"last_active"`
}

func (user *User) UpdateToken(token string) {
	user.Token = token
	db.Save(&user)
}
