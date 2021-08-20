package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Msg    interface{} `json:"msg"`
}

const (
	ERROR   = false
	SUCCESS = true
)

func result(status bool, data interface{}, msg interface{}, c *gin.Context) {
	// 返回json数据
	c.JSON(http.StatusOK, Response{
		Status: status,
		Msg:  msg,
		Data: data,


	})
}

func SuccessWithData(data interface{},message string,c *gin.Context){
	result(SUCCESS,data,message,c)
}

func SuccessNoData(message string,c *gin.Context){
	result(SUCCESS,map[string]interface{}{},message,c)
}

//只返回错误信息
func Failed(message string,c *gin.Context){
	result(ERROR,map[string]interface{}{},message,c)
}

func FailedWithValid(message interface{},c *gin.Context){
	result(ERROR,map[string]interface{}{},message,c)
}