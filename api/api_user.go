package api

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"netdisk/call"
	"netdisk/global"
	"netdisk/middleware"
	"netdisk/model"
	"netdisk/model/requests"
	"netdisk/model/response"
	"netdisk/utils"
)

// @Tags 用户相关
// @Summary 用户注册账号
// @Produce application/json
// @Accept application/x-www-form-urlencoded
// @Param username body string  true "用户名 "
// @Param password body string  true "用户密码 "
// @Success 200 {object} response.OKWithoutData
// @Failure 500 {object} response.SystemFailed
// @Router /user/register [post]
func Register(c *gin.Context) {
	var R requests.Register
	if err:=utils.Validator(c,&R,utils.BindForm);err!=nil {
		return
	}

	user := model.User{Username: R.Username, Password: R.Password}
	_, Error, _ := call.Register(call.User(user)) //返回的Error为string类型
	if Error != "" {
		global.SugaredLogger.Error("注册失败", zap.Any("err:", errors.New(Error)))
		response.Failed("注册失败:"+Error, c)
		return
	}
	response.SuccessNoData("注册成功",c)
}


// @Tags 用户相关
// @Summary 用户密码登录
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param username body string  true "用户名"
// @Param password body string  true "用户密码"
// @Success 200 {object} response.LoginResp
// @Failure 500 {object} response.SystemFailed
// @Router /user/login [post]
func Login(c *gin.Context){
	var R requests.Register
	if err:=utils.Validator(c,&R,utils.BindForm);err!=nil{
		return
	}
	user,Error,_:=call.VerifyPassword(R.Username,R.Password)
	if Error!=""{
		response.Failed(Error,c)
		return
	}
	//验证通过颁发token
	TokenNext(model.User(user),c)
}

func TokenNext(user model.User,c *gin.Context){
	claims:= model.CustomClaims{
		ID:       user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(48 * time.Hour).Unix(), // 过期时间 2天
			Issuer:    "xzh",                                 // 签名的发行者
		},
	}
	token,err:=middleware.GenToken(claims)
	if err!=nil{
		global.SugaredLogger.Error("生成token失败",zap.Any("err:",err))
		response.Failed("获取token失败",c)
		return
	}
	response.SuccessWithData(response.LoginRespData{
		Token: token,
		ExpTime: time.Unix(claims.StandardClaims.ExpiresAt,0).Format("2006-01-02 15:04:05"),
	},"登录成功",c)
}
