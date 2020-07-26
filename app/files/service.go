package files

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadFile : 上传文件
func (s *conf.Service) UploadFile(ctx *gin.Context) {
	var req FileUploadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to bind json",
			"data": "",
		})
		return
	}

	resp, err := FileUpload(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to bind json",
			"data": "",
		})
		return
	}

	ctx.JSON(200, resp)
}
