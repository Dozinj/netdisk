package middleware

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"netdisk/global"
	"netdisk/model"
	"netdisk/model/response"
)


func JwtAuth()gin.HandlerFunc{
	return func(ctx *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		authHeader:=ctx.Request.Header.Get("Authorization")
		if authHeader==""{
			response.Failed("没有用户权限",ctx)
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Failed("请求头中auth格式有误",ctx)
			ctx.Abort()
			return
		}

		// parts[1]是获取到的tokenString
		claims,err:=ParseToken(parts[1])
		if err!=nil{
			if err==TokenExpired{
				response.Failed("用户令牌过期",ctx)
				ctx.Abort()
				return
			}
			response.Failed(err.Error(), ctx)
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

func GetJwtKey()[]byte{
	return []byte(global.Config.Section("jwt").Key("jwtkey").String())
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

//生成token
func GenToken(claims model.CustomClaims)(string,error) {
	// 使用指定的签名方法创建签名对象
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(GetJwtKey()) //这里key为string类型会报错
}

//解析token
func ParseToken(tokenString string)(*model.CustomClaims,error){
	//解析token
	token,err:=jwt.ParseWithClaims(tokenString,&model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJwtKey(),nil
	})

	if err != nil {
		global.SugaredLogger.Error("token解析错误",zap.Any("err:",err))
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}

