package model

type User struct {
	ID        int    `json:"id"`                         //主键id
	Username  string `json:"username" example:"xzh"`     // 用户登录名
	Password  string `json:"password"  example:"123456"` // 用户登录密码
	Salt      string `json:"salt"`                       //用户密码加密撒盐值
	HeaderImg string `json:"header_img"`                 // 用户头像
}
