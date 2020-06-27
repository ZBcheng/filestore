package drivers

import (
	"testing"
)

func TestMysqlConn(t *testing.T) {
	t.Log("conn...")
}

func TestQuery(t *testing.T) {
	db := DBConn()
	rows, err := db.Query("SELECT id, user_name, user_pwd FROM tbl_user")
	if err != nil {
		t.Error(err)
	}

	for rows.Next() {
		var id int
		var user_name string
		var user_pwd string
		if err := rows.Scan(&id, &user_name, &user_pwd); err != nil {
			t.Error(err)
		}
		t.Log(id, user_name, user_pwd)
	}
}
