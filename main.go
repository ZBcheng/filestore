package main

import (
	"fmt"
	"moviesite-filestore/handler"
	"moviesite-filestore/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(util.Cors())
	router.POST("/file/upload", handler.UploadHandler)
	router.POST("/file/mpupload", handler.MultipartUploadHandler)
	router.POST("/events", events)
	router.GET("/file/mpupload/uppart", handler.UploadPartHandler)
	// router.POST("/file/mpupload/complete", handler.CompleteUploadHandler)
	// router.GET("file/mpupload/uppart", handler.UploadPartHandler)
	// router.GET("/file/init", handler.InitialMultipartUploadHandler)
	router.Run(":7000")
	// fmt.Println(os.Open("/Users/zhangbicheng/Desktop/valley.png"))
}

func events(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	fmt.Println(c.Request.Body)
	fmt.Println("buf: " + string(buf[0:n]))
	resp := map[string]string{"hello": "world"}
	c.JSON(http.StatusOK, resp)
	/*post_gwid := c.PostForm("name")
	fmt.Println(post_gwid)*/
}
