package meta

import (
	"testing"

	"github.com/zbcheng/filestore/models"
	"github.com/zbcheng/filestore/util"
)

func TestUpdateDB(t *testing.T) {
	fileNameBytes := util.Sha1([]byte("hello.txt"))
	fMeta := models.FileMeta{
		FileHash: fileNameBytes,
		FileName: "hello.txt",
		FileSize: 64,
		Location: "/Users/zhangbicheng/PycharmProjects",
	}
	if suc := UpdateFileMetaDB(fMeta); suc {
		t.Log("Update success")
	} else {
		t.Error("Update Failed")
	}
}

func TestFileExists(t *testing.T) {
	fileNameBytes := util.Sha1([]byte("hello.txt"))
	fMeta := models.FileMeta{
		FileHash: fileNameBytes,
		FileName: "hello.txt",
		FileSize: 64,
		Location: "/Users/zhangbicheng/PycharmProjects",
	}
	exists, err := FileExists(fMeta)
	if err != nil {
		t.Error(err)
	}
	t.Log(exists)
}
