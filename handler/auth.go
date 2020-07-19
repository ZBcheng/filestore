package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// AuthHandler : 鉴权
func AuthHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println(username)
	fmt.Println(password)
}
