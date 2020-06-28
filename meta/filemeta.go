package meta

import (
	"database/sql"
	"fmt"

	drivers "github.com/zbcheng/filestore/drivers/mysql"
	"github.com/zbcheng/filestore/models"
	repo "github.com/zbcheng/filestore/repository"
)

// FileMeta : file struct
type FileMeta struct {
	FileHash string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta
var db *sql.DB

func init() {
	fileMetas = make(map[string]FileMeta)
	db = drivers.DBConn()
}

// UpdateFileMeta : add or update a file
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileHash] = fmeta
}

// UpdateFileMetaDB: 新增/更新文件元信息到Mysql
func UpdateFileMetaDB(fmeta models.FileMeta) bool {
	return repo.OnFileUploadFinished(
		fmeta.FileHash, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// GetFileMeta : get a file
func GetFileMeta(fileHash string) (FileMeta, error) {
	fileMeta := FileMeta{}
	querySQL := fmt.Sprintf("SELECT filename, filesize, location, uploadtime FROM tbl_file where filehash='%s'", fileHash)
	rows, err := db.Query(querySQL)

	for rows.Next() {
		err = rows.Scan(&fileMeta.FileName, &fileMeta.FileSize, &fileMeta.Location, &fileMeta.UploadAt)
		if err != nil {
			return FileMeta{}, err
		}
	}

	return fileMeta, nil
}

// RemoveFileMeta : remove file info from db
func RemoveFileMeta(fileHash string) error {
	stmt, err := db.Prepare("DELETE FROM tbl_file where filehash=($1)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(fileHash); err != nil {
		return err
	}

	return nil

}

// FileExists : 判断文件是否在db中存在
func FileExists(f models.FileMeta) (exists bool, err error) {
	sql := fmt.Sprintf("SELECT file_size FROM tbl_file WHERE file_sha1='%s'", f.FileHash)
	rows, err := db.Query(sql)

	if err != nil {
		return false, err
	}

	var filesize int64

	for rows.Next() {
		if err := rows.Scan(&filesize); err != nil {
			return false, err
		}

		if filesize == f.FileSize {
			return true, nil
		}
	}

	return false, nil
}
