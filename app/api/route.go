package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zbcheng/filestore/util"
)

func RegisterRouter() *gin.Engine {
	// service := &Service{}

	router := gin.Default()
	router.Use(util.Cors())
	// router.POST("/file/upload", handler.UploadHandler)
	// router.POST("/file/upload", files.FileUpload)

	// router.POST("/user/signin", handler.Signin)
	// router.POST("/user/signup", handler.Signup)
	// router.GET("/user/info", handler.UserInfo)

	return router
}
