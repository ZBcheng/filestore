package repository

import (
	"time"

	db "github.com/zbcheng/filestore/app/drivers/mysql"
	"github.com/zbcheng/filestore/app/models"
)

func StoreFileMeta(fMeta models.FileMeta) {
	file := models.FileMeta{
		FileHash: fMeta.FileHash,
		FileName: fMeta.FileName,
		FileSize: fMeta.FileSize,
		Location: fMeta.Location,
		UploadAt: string(time.Now().Format("2006-01-02 15:04:05")),
	}

	db.DBConn().Create(&file)

}
