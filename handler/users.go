package handler

import (
	"net/http"

	"github.com/arstd/log"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	drivers "github.com/zbcheng/filestore/drivers/mysql"
	"github.com/zbcheng/filestore/models"
	repo "github.com/zbcheng/filestore/repository"
	"github.com/zbcheng/filestore/util"
)

var db *gorm.DB

func init() {
	db = drivers.DBConn()
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// Signin : 用户登录
func Signin(c *gin.Context) {
	loginForm := LoginForm{}
	if err := c.BindJSON(&loginForm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Internal Server Error",
			"err":  1,
			"data": "",
		})
		return
	}

	username := loginForm.Username
	password, err := util.EncryptPwd(loginForm.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Internal Server Error",
			"err":  1,
			"data": "",
		})
		log.Info("Failed to auth password")
		return
	}

	msg, suc := repo.AuthUser(username, password, repo.GenToken(username))

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

func Signup(c *gin.Context) {
	form := &SignupForm{}
	if err := c.BindJSON(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to bind json",
			"err":  1,
			"data": "",
		})
		return
	}

	encPwd, err := util.EncryptPwd(form.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Faile to encrypt password",
			"err":  1,
			"data": "",
		})
		log.Info("Failed to encrypt password")
		return
	}

	user := models.User{
		Username: form.Username,
		Password: encPwd,
		Email:    form.Email,
		Phone:    form.Phone,
		Avatar:   form.Avatar,
		Token:    repo.GenToken(form.Username),
	}

	if err := repo.CreateUser(&user); err != nil {
		if err.Error() ==
			"Error 1062: Duplicate entry '"+user.Username+"' for key 'username'" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "Username '" + user.Username + "' has already exist!",
				"err":  1,
				"data": "",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Internal Server Error",
			"err":  1,
			"data": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "OK",
		"err":  0,
		"data": user,
	})

}

// UserInfo : 获取用户信息
func UserInfo(c *gin.Context) {
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
