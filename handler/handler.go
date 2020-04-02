package handler

import (
	"fmt"
	"io"
	"moviesite-filestore/meta"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler : 上传文件接口
func UploadHandler(c *gin.Context) {
	fmt.Println("requesting upload handler")
	file, fHead, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("form file failed")
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": "form file failed",
		})
		return
	}

	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName: fHead.Filename,
		Location: "/Users/zhangbicheng/Desktop/" + fHead.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": "create file failed",
		})
		return
	}

	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": "write file error",
		})
		return
	}

	newFile.Seek(0, 0)

	c.JSON(http.StatusOK, gin.H{
		"data": "upload complete",
	})
}

// FileQueryHandler : query file by filehash
// func FileQueryHandler(c *gin.Context) {
// 	filehash := c.Query("filehash")[0]
// 	fileMetas := meta.GetFileMeta()
// 	data, err := json.Marshal((fileMetas))
// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"data": err,
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": data,
// 	})
// }
