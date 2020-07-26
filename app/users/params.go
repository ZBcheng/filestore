package users

type UserSignupReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type UserSignupResp struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

type UserSigninReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type UserSigninResp struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

type UserInfoReq struct {
	ID int `json:"id"`
}

type UserInfoResp struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role"`
}
