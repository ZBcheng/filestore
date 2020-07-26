package users

import (
	"net/http"

	"github.com/arstd/log"
	"github.com/gin-gonic/gin"
	"github.com/zbcheng/filestore/app/models"
	"github.com/zbcheng/filestore/conf"
	"github.com/zbcheng/filestore/util"
)

// Signin : 用户登录
func Signin(req *UserSigninReq) (res *Resp, err error) {
	res = new(Resp)

	username := req.Username
	password, err := util.EncryptPwd(req.Password)
	if err != nil {
		res.Code = 500
		res.Msg = "Encrypt password error"
		return
	}

	token := repo.GenToken(username)
	id, err := repo.AuthUser(username, password, token)

	if err != nil {
		res.Code = 400
		res.Msg = "Wrong username or password"
		return
	}

	if id == 0 {
		res.Code = 400
		res.Msg = "Auth failed"
		log.Debug(username, password, req.Password)
		return
	}

	res.Code = conf.SucRespCode
	res.Msg = conf.SucRespMsg
	res.Data = &UserSigninResp{
		ID:    id,
		Token: token,
	}

	return res, nil

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
