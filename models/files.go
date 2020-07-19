package models

import (
	"github.com/jinzhu/gorm"
	drivers "github.com/zbcheng/filestore/drivers/mysql"
)

func init() {

	db := drivers.DBConn()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&FileMeta{})

}

type FileMeta struct {
	gorm.Model
	FileHash string `gorm:"filehash" json:"filehash"`
	FileName string `gorm:"filename" json:"filename"`
	FileSize int64  `gorm:"filesize" json:"filesize"`
	Location string `gorm:"location" json:"location"`
	UploadAt string `gorm:"upload_at" json:"upload_at"`
}
