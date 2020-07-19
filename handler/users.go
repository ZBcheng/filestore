package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	drivers "github.com/zbcheng/filestore/drivers/mysql"
	"github.com/zbcheng/filestore/models"
	repo "github.com/zbcheng/filestore/repository"
)

var db *gorm.DB

func init() {
	db = drivers.DBConn()
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SigninHandler(c *gin.Context) {
	loginForm := LoginForm{}
	if err := c.BindJSON(&loginForm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	username := loginForm.Username
	pasword := loginForm.Password

	user := models.User{}
	db.Where("username = ?", username).First(&user)

	msg, suc := repo.AuthUser(username, pasword, repo.GenToken(username))

	if !suc {
		c.JSON(http.StatusForbidden, gin.H{
			"msg":  msg,
			"err":  1,
			"data": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"err":  0,
		"data": "OK",
	})

}

func UserInfoHandler(c *gin.Context) {
	username := c.Query("username")
	token := c.Query("token")

	if repo.AuthToken(username, token) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "OK",
			"data": "",
		})
		return
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"msg":  "token doesn't match username or token expired",
			"data": "",
		})
		return
	}
}
