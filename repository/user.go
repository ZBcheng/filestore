package repository

import (
	"time"

	"github.com/arstd/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/zbcheng/filestore/conf"
	drivers "github.com/zbcheng/filestore/drivers/mysql"
	"github.com/zbcheng/filestore/models"
)

var db *gorm.DB

func init() {
	db = drivers.DBConn()
}

// UserSignup : 用户注册
func UserSignup(username, password, email, phone string) bool {

	if username == "" || password == "" {
		return false
	}

	signupAt := time.Now()
	lastActiveAt := signupAt
	user := models.User{
		Username:   username,
		Password:   password,
		Email:      email,
		Phone:      phone,
		SignupAt:   signupAt,
		LastActive: lastActiveAt,
		Status:     0,
		Token:      GenToken(username),
	}

	db.Create(&user)

	return true
}

// GenToken : 生成token
func GenToken(username string) string {
	secretKey := conf.Load().Secret.SecretKey
	// secretKey := "always blue"

	claims := make(jwt.MapClaims)
	claims["username"] = username

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}

func AuthToken(username, token string) bool {
	var tk *jwt.Token
	var secretKey = conf.Load().Secret.SecretKey
	tk, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Error("Failed to auth token: ", err)
		return false
	}

	claims := tk.Claims.(jwt.MapClaims)

	if username != claims["username"] {
		log.Debug("Username doesn't mastch!")
		return false
	}

	curTime := time.Now().Unix()

	if curTime > int64(claims["exp"].(float64)) {
		log.Debug("Token expired!")
		return false
	}

	return true
}

func AuthUser(username, password, token string) (message string, success bool) {
	user := &models.User{}
	db.Where("name = ?", username).First(&user)

	if password != user.Password {
		return "wrong username or password!", false
	}

	if tokenValid := AuthToken(username, token); !tokenValid {
		return "invalid token!", false
	}

	return "ok", true
}
