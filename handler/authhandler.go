package handler


// AuthHandler : 鉴权
func AuthHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println("username")
	fmt.Println("passworld")
}