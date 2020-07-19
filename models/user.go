package models

import (
	"github.com/jinzhu/gorm"
	drivers "github.com/zbcheng/filestore/drivers/mysql"
)

var db *gorm.DB

func init() {
	db = drivers.DBConn()
}

type User struct {
	gorm.Model
	Username string `gorm:"username; type:varchar(45); not null;unique" json:"username"` // 用户名
	Password string `gorm:"password; type:varchar(255)" json:"password"`                 // 密码
	Email    string `gorm:"email" json:"email"`                                          // 邮箱
	Phone    string `gorm:"phone; type:varchar(20)" json:"phone"`                        // 手机
	Avatar   string `gorm:"avatar" json:"avatar"`                                        // 头像
	Status   int    `gorm:"status" json:"status"`
	Token    string `gorm:"token" json:"token"`
}

func (user *User) UpdateToken(token string) {
	user.Token = token
	db.Save(&user)
}
