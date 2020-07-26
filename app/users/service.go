package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) UserSignin(ctx *gin.Context) {
	var req handler.Request
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to bind json",
			"data": "",
		})
	}

	resp, err := handler.Signin(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": resp.Code,
			"msg":  resp.Error,
			"data": "",
		})
	}

	ctx.JSON(http.StatusOK, resp)
}
