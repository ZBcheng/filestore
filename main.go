package main

import (
	"github.com/zbcheng/filestore/handler"
	"github.com/zbcheng/filestore/util"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(util.Cors())
	// router.POST("/file/upload", handler.UploadHandler)
	router.POST("/file/upload", handler.UploadHandler)

	router.POST("/user/signin", handler.SigninHandler)
	router.GET("/user/info", handler.UserInfoHandler)
	router.Run(":7000")
}
