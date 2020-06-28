package repository

import (
	"fmt"
)

func OnFileUploadFinished(filehash string, filename string,
	filesize int64, fileaddr string) bool {
	stmt, err := db.Prepare("INSERT IGNORE INTO tbl_file (`file_sha1`, `file_name`, `file_size`," +
		"`file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err)
		return false
	}

	rowsAffected, err := ret.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if rowsAffected <= 0 {
		fmt.Printf("File with hash:%s has been uploaded before\n", filehash)
	}

	return true
}
