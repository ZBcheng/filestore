package files

type FileMeta struct {
	FileHash string `gorm:"filehash" json:"filehash"`
	FileName string `gorm:"filename" json:"filename"`
	FileSize int64  `gorm:"filesize" json:"filesize"`
	Location string `gorm:"location" json:"location"`
	UploadAt string `gorm:"upload_at" json:"upload_at"`
}
