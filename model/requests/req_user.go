package requests

// User register structure
type Register struct {
	Username string `json:"username" form:"username" binding:"required" ` //用户名
	Password string `json:"password" form:"password" binding:"required" ` //用户密码
}


