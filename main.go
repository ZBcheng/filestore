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

	router.POST("/user/signin", handler.Signin)
	router.POST("/user/signup", handler.Signup)
	router.GET("/user/info", handler.UserInfo)
	router.Run(":7000")
}
