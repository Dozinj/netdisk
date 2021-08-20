package api

import (
	"github.com/gin-gonic/gin"

	"netdisk/model/requests"
	"netdisk/model/response"
	"netdisk/service"
	"netdisk/utils"
)

//@Tags 分享相关
//@Summary 获取分享的二维码
//@Description  category 1为video,2为image
//@accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param category body string true "文件类型  category值1代表video,2代表image"
//@Param fileid body string  true "文件id"
//@Success 200 {object} response.RespGQrcode
//@Failure 500 {object} response.SystemFailed
//@Router /share/qrcode [POST]
func GenerateQrcode(c *gin.Context){
	var Gq requests.GQrcode
	err:=utils.Validator(c,&Gq,utils.BindForm)
	if err!=nil{
		return
	}
	if !(Gq.Category=="1"||Gq.Category=="2"){
		response.Failed("category 字段只能是1或2",c)
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}

	fs:=service.GetfileService()
	url,err:=fs.Generateqrcode(&Gq,username)
	if err!=nil{
		response.Failed(err.Error(),c)
		return
	}
	response.SuccessWithData(response.RespGQrcodeData{Url: url},"生成分享文件二维码成功",c)
}


//@Tags 分享相关
//@Summary 解析二维码
//@Description 返回要分享文件相关信息 category 1为video,2为image
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param category query string true "文件类型  category值1代表video,2代表image"
//@Param shareuser query string true "分享者用户名"
//@Param fileid query string  true "文件id"
//@Success 200 {object} response.RespShareFile
//@Failure 500 {object} response.SystemFailed
//@Router /share/qrcode [GET]
func AnalyzeQrcode(c *gin.Context) {
	var Aq requests.AQrcode
	err := utils.Validator(c, &Aq, utils.BindQuery)
	if err != nil {
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}
	fs:=service.GetfileService()
	shareFile, err := fs.GetFileInfoByQrcode(&Aq, username)
	if err != nil {
		response.Failed(err.Error(), c)
		return
	}
	response.SuccessWithData(shareFile, "获取分享文件信息成功", c)
}


//@Tags 分享相关
//@Summary 获取分享链接
//@Description 文件提取码可以不设置,若有提取码长度限制为4
//@Accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param category body string true "文件类型  category值1代表video,2代表image"
//@Param fileid body string  true "文件id"
//@Param extraction_code body string  false "文件提取码"
//@Success 200 {object} response.RespShareLink
//@Failure 500 {object} response.SystemFailed
//@Router /share/link  [POST]
func GenerateSharingLink(c *gin.Context) {
	var Gl requests.GLink
	err:=utils.Validator(c,&Gl,utils.BindForm)
	if err!=nil{
		return
	}

	if !(Gl.Category=="1"||Gl.Category=="2"){
		response.Failed("category 字段只能是1或2",c)
		return
	}
	if Gl.ExtractionCode!=""&&len(Gl.ExtractionCode)!=4{
		response.Failed("提取码长度只能是4位",c)
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}

	fs:=service.GetfileService()
	url, err := fs.GenerateSharingLink(&Gl,username)
	if err != nil {
		response.Failed(err.Error(), c)
		return
	}

	response.SuccessWithData(response.RespShareLinkData{
		ShareLink:      url,
		ExtractionCode: Gl.ExtractionCode},
		"分享链接生成成功", c)
}

//@Tags 分享相关
//@Summary 解析链接
//@Description 返回要分享文件相关信息
//@Accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param link path string  true "分享链接"
//@Param extraction_code body string false "提取码"
//@Success 200 {object} response.RespShareFile
//@Failure 500 {object} response.SystemFailed
//@Router /s/:link  [POST]
func AnalyzeShareLink(c *gin.Context) {
	link := c.Param("link")
	extractionCode := c.PostForm("extraction_code")
	if link == "" {
		response.Failed("链接无效", c)
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}
	fs:=service.GetfileService()
	shareFile, err := fs.GetFileInfoByLink(link, extractionCode,username)
	if err != nil {
		response.Failed(err.Error(), c)
		return
	}
	response.SuccessWithData(shareFile, "获取分享文件信息成功", c)
}


