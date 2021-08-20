package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"netdisk/global"
	"netdisk/model"
	"netdisk/model/response"
)


func GetTokenInfo(c *gin.Context)string{
	//获取jwt解析的用户信息
	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*model.CustomClaims)
	if !ok {
		global.SugaredLogger.Info("jwt中间件未提供请求信息")
		return ""
	}
	return customClaims.Username
}

const (
	start=iota
	NormalUpload
	AppendUpload
)
var validErr =errors.New("参数效验失败")
var (
	BindQuery="query"
	BindForm="form"
)
func Validator(c *gin.Context,obj interface{},BindType string)(err error){
	switch BindType{
	case BindQuery:
		err=c.ShouldBindQuery(obj)
	case BindForm:
		err=c.ShouldBind(obj)
	default:
		response.Failed(validErr.Error(),c)
		return validErr
	}

	if  err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			global.SugaredLogger.Error("参数绑定失败",zap.Any("err:",err))
			response.Failed(validErr.Error(),c)
			return validErr
		}
		// validator.ValidationErrors类型错误则进行翻译
		response.FailedWithValid(removeTopStruct(errs.Translate(global.Trans)),c)
		return validErr
	}
	return nil
}

//去掉validator效验结构体前缀  "SignUpParam.Age"
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}


