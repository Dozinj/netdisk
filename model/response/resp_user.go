package response

type OKWithoutData struct {
	status bool        `example:"true"`
	Data   interface{} `json:"data" example:""`
	Msg    string      `example:"xxx成功"`
}

type SystemFailed struct {
	Status bool        `example:"false"`
	Data   interface{} `json:"data" example:""`
	Msg    string      `example:"status_internal_serverError"`
}

type LoginResp struct {
	status bool          `example:"true"`
	Data   LoginRespData `json:"data"`
	Msg    string        `example:"登录成功"`
}
type LoginRespData struct {
	Token    string `json:"token"`    //token令牌
	ExpTime  string `json:"过期时间"` //token 过期时间
}