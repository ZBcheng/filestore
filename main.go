package main

import (
	"filestore/handler"
	"filestore/util"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(util.Cors())
	// router.POST("/file/upload", handler.UploadHandler)
	router.POST("/file/mpupload", handler.MultipartUploadHandler)
	router.Run(":7000")
}
