package repository

import (
	"time"

	drivers "github.com/zbcheng/filestore/drivers/mysql"
	"github.com/zbcheng/filestore/models"
)

func StoreFileMeta(fMeta models.FileMeta) {
	db := drivers.DBConn()
	file := models.FileMeta{
		FileHash: fMeta.FileHash,
		FileName: fMeta.FileName,
		FileSize: fMeta.FileSize,
		Location: fMeta.Location,
		UploadAt: string(time.Now().Format("2006-01-02 15:04:05")),
	}

	db.Create(&file)

}
