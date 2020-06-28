package repository

import (
	"database/sql"
	"errors"
	"time"

	drivers "github.com/zbcheng/filestore/drivers/mysql"
)

var db *sql.DB

func init() {
	db = drivers.DBConn()
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// UserSingup : 用户注册
func UserSignup(username, password string) (success bool, err error) {
	stmt, err := db.Prepare("INSERT IGNORE INTO tbl_user (`user_name`, `user_pwd`) values(?, ?)")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		return false, err
	}

	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true, nil
	}

	return false, errors.New("Unknown Error")

}

func UserSignin(username, password string) (success bool, err error) {
	var pwd string

	stmt, err := db.Prepare("SELECT user_pwd FROM tbl_user WHERE `user_name`=? limit 1")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return false, nil
	} else if rows == nil {
		return false, errors.New("username not found")
	}

	for rows.Next() {
		if err = rows.Scan(&pwd); err != nil {
			return false, err
		}
	}

	if password == pwd {
		if err = updateLastActive(username); err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, errors.New("wrong password")
	}

}

func UpdateToken(username, token string) (success bool, err error) {
	stmt, err := db.Prepare("REPLACE INTO tbl_user_token (`username`, `user_token`) values(?,?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, token)
	if err != nil {
		return false, err
	}

	rowsAffected, err := ret.RowsAffected()
	if err == nil && rowsAffected > 0 {
		return true, nil
	}
	return false, err
}

// UpdateLastActive : 更新登录时间
func updateLastActive(username string) (err error) {
	lastActiveAt := time.Now().Format("2006-01-02 15:04:05")
	stmt, err := db.Prepare("UPDATE tbl_user SET last_active=? WHERE user_name=?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	ret, err := stmt.Exec(lastActiveAt, username)
	if err != nil {
		return err
	}

	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return nil
	}

	return errors.New("Falied to update last_active")

}
