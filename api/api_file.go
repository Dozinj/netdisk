package api

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"netdisk/model/requests"
	"netdisk/model/response"
	"netdisk/service"
	"netdisk/utils"
)

var JwtErr="用户权限验证失败"

// @Tags 文件相关
// @Summary 视频上传
// @Description 若是文件大小大于50MB,则采用分片上传
// @Accept multipart/form-data
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param video formData file true "上传视频"
// @Success 200 {object} response.RespUploadVideo
// @Failure 500 {object} response.SystemFailed
// @Router /file/video  [post]
func UploadVideo(c *gin.Context) {
	videoHeader, err := c.FormFile("video")
	if err != nil {
		response.Failed("video字段不能为空", c)
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed(JwtErr, c)
		return
	}

	//上传视频和封面
	fs:=service.GetfileService()
	file, err := fs.UploadVideo(videoHeader,username)
	if err != nil {
		response.Failed(err.Error(), c)
		return
	}

	res := utils.FloatRound(float64(file.Size)/float64(1<<20), 3)
	response.SuccessWithData(response.RespUploadVideoData{
		VideoId:        file.ID,
		VideoName:      file.FileName,
		Size:           strconv.FormatFloat(res, 'f', -1, 64) + "MB",
	}, "上传视频成功", c)
}


//@Tags 文件相关
//@Summary 图片上传
//@Description
//@accept multipart/form-data
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param image formData file true "上传图片"
//@Success 200 {object} response.RespUploadImage
// @Failure 500 {object} response.SystemFailed
//@Router /file/image  [post]
func UploadImage(c *gin.Context) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.Failed("image 字段不能为空", c)
		return
	}

	//获取jwt解析的用户信息
	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户权限验证失败", c)
		return
	}
	//获取服务
	fs := service.GetfileService()
	image, err := fs.UploadImage(fileHeader, username)
	if err != nil {
		response.Failed(err.Error(), c)
		return
	}

	//文件大小 单位MB
	res := utils.FloatRound(float64(fileHeader.Size)/float64(1<<20), 3)
	response.SuccessWithData(response.RespUploadImageData{
		ImageId:        image.ID,
		ImageName:      image.FileName,
		Size:           strconv.FormatFloat(res, 'f', -1, 64) + "MB",
	}, "上传图片成功", c)
}



//@Tags 文件相关
//@Summary 保存网盘
//@Description 将别人分享的文件上传至自己网盘
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param category query string true "文件类型  category值 1代表video,2代表image"
//@Param fileid query string ture "文件id"
//@Success 200 {object} response.OKWithoutData
// @Failure 500 {object} response.SystemFailed
//@Router /disk/save [GET]
func Save(c *gin.Context) {
	var fc requests.FC
	err := utils.Validator(c, &fc, utils.BindQuery)
	if err != nil {
		return
	}

	if !(fc.Category=="1"||fc.Category=="2"){
		response.Failed("category 字段只能是1或2",c)
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}

	fs:=service.GetfileService()
	if err = fs.DiskSave(&fc, username); err != nil {
		response.Failed(err.Error(), c)
		return
	}

	response.SuccessWithData(nil, "文件保存网盘成功", c)
}


//@Tags 文件相关
//@Summary 保存本地
//@Description 将别人分享的文件下载到本地
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param category query string true "文件类型  category 1代表video,2代表image"
//@Param fileid query string ture "文件id"
//@Success 200 {object} response.RespLocalSave
//@Failure 500 {object} response.SystemFailed
//@Router /local/save [GET]
func Download(c *gin.Context){
	var fc requests.FC
	err := utils.Validator(c, &fc, utils.BindQuery)
	if err != nil {
		return
	}

	if !(fc.Category=="1"||fc.Category=="2"){
		response.Failed("category 字段只能是1或2",c)
		return
	}
	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}
	fs:=service.GetfileService()
	url,err:=fs.LocalSave(&fc,username)
	if err!=nil{
		response.Failed(err.Error(),c)
		return
	}
	response.SuccessWithData(response.RespLocalSaveData{Url: url},"获取下载地址成功",c)
}


//@Tags 文件相关
//@Summary 视频目录
//@Description 获取当前目录的所有文件和文件夹
//@accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param path body string  true "文件目录 path需以`/`开头和结尾"
//@Success 200 {object} response.RespList
//@Failure 500 {object} response.SystemFailed
//@Router /list/video [POST]
func ListVideos(c *gin.Context) {
	path:=c.PostForm("path")
	if path==""{
		response.Failed("path 字段不允许为空",c)
		return
	}

	paths:= strings.Split(path,"/")
	if paths[len(paths)-1]!=""{
		response.Failed("file_path必须以`/`结尾",c)
		return
	}else if paths[0]!=""{
		response.Failed("file_path必须以`/`开头",c)
		return
	}


	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}
	fs:=service.GetfileService()
	datas, err := fs.ListVideo(path,username)
	if err != nil {
		response.Failed(err.Error(), c)
		return
	}

	response.SuccessWithData(*datas, "获取网盘信息成功", c)
}


//@Tags 文件相关
//@Summary 图片目录
//@Description 获取当前目录的所有文件和文件夹
//@accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param path body string  true "文件目录 path需以`/`开头和结尾"
//@Success 200 {object} response.RespList
//@Failure 500 {object} response.SystemFailed
//@Router /list/image [POST]
func ListImages(c *gin.Context) {
	path:=c.PostForm("path")
	if path==""{
		response.Failed("path 字段不允许为空",c)
		return
	}

	paths:= strings.Split(path,"/")
	if paths[len(paths)-1]!=""{
		response.Failed("file_path必须以`/`结尾",c)
		return
	}else if paths[0]!=""{
		response.Failed("file_path必须以`/`开头",c)
		return
	}


	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}
	fs:=service.GetfileService()
	datas, err := fs.ListImage(path,username)
	if err != nil {
		response.Failed(err.Error(), c)
		return
	}
	response.SuccessWithData(*datas, "获取网盘信息成功", c)
}


//@Tags 文件相关
//@Summary 更改视频路径
//@Description
//@Accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param fileid body string true "文件id"
//@Param file_path body string true "文件新路径 path需以`/`开头和结尾"
//@Success 200 {object} response.OKWithoutData
//@Failure 500 {object} response.SystemFailed
//@Router /file/video/path [PUT]
func ChangeVideoPath(c *gin.Context){
	var Cp requests.ChangePath
	err:=utils.Validator(c,&Cp,utils.BindForm)
	if err!=nil{
		return
	}

	paths:= strings.Split(Cp.FilePath,"/")
	if paths[len(paths)-1]!=""{
		response.Failed("file_path必须以`/`结尾",c)
		return
	}
	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}

	fs:=service.GetfileService()
	if err:=fs.ChangeVideoPath(Cp.Fileid,Cp.FilePath,username);err!=nil{
		response.Failed(err.Error(),c)
		return
	}

	response.SuccessNoData("修改视频路径成功",c)
}

//@Tags 文件相关
//@Summary 更改图片路径
//@Description
//@Accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param fileid body string true "文件id"
//@Param file_path body string true "文件新路径 path需以`/`开头和结尾"
//@Success 200 {object} response.OKWithoutData
//@Failure 500 {object} response.SystemFailed
//@Router /file/image/path [PUT]
func ChangeImagePath(c *gin.Context){
	var Cp requests.ChangePath
	err:=utils.Validator(c,&Cp,utils.BindForm)
	if err!=nil{
		return
	}

	paths:= strings.Split(Cp.FilePath,"/")
	if paths[len(paths)-1]!=""{
		response.Failed("file_path必须以`/`结尾",c)
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}

	fs:=service.GetfileService()
	if err:=fs.ChangeImagePath(Cp.Fileid,Cp.FilePath,username);err!=nil{
		response.Failed(err.Error(),c)
		return
	}

	response.SuccessNoData("修改图片路径成功",c)
}


//@Tags 文件相关
//@Summary 修改文件名
//@Description  category 1为video,2为image,暂不支持修改文件格式
//@accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param category body string true "文件类型  category 1代表video,2代表image"
//@Param fileid body string  true "文件id"
//@Param filename body string  true "新文件名"
//@Success 200 {object} response.OKWithoutData
//@Failure 500 {object} response.SystemFailed
//@Router /file/filename [PUT]
func ModifyFileName(c *gin.Context){
	var Cf requests.ChangeFilename
	err:=utils.Validator(c,&Cf,utils.BindForm)
	if err!=nil{
		return
	}
	if !(Cf.Category=="1"||Cf.Category=="2"){
		response.Failed("category 字段只能是1或2",c)
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}

	fs:=service.GetfileService()
	if err:=fs.ChangeFile(Cf,Cf.Fileid,Cf.Category,username);err!=nil{
		response.Failed(err.Error(),c)
		return
	}

	response.SuccessNoData("文件名修改成功",c)
}


//@Tags 文件相关
//@Summary 修改文件权限
//@Description
//@accept application/x-www-form-urlencoded
//@Produce application/json
//@Param Authorization header string true "Bearer token"
//@Param category body string true "文件类型  category 1代表video,2代表image"
//@Param fileid body string  true "文件id"
//@Param authority body string  true "新文件权限 authority:8.所有人可下载 9.仅获得链接的人可下载 10.仅自己见"
//@Success 200 {object} response.OKWithoutData
//@Failure 500 {object} response.SystemFailed
//@Router /file/authority [PUT]
func ModifyAuthority(c *gin.Context) {
	var Ca requests.ChangeAuthority
	err := utils.Validator(c, &Ca, utils.BindForm)
	if err != nil {
		return
	}

	if !(Ca.Category=="1"||Ca.Category=="2"){
		response.Failed("category 字段只能是1或2",c)
		return
	}
	authority:=Ca.Authority
	if !(authority=="8"||authority=="9"||authority=="10"){
		response.Failed("authority 字段只能是8或9或10",c)
		return
	}

	username := utils.GetTokenInfo(c)
	if username == "" {
		response.Failed("用户信息获取失败", c)
		return
	}

	fs := service.GetfileService()
	if err := fs.ChangeFile(Ca, Ca.Fileid, Ca.Category, username); err != nil {
		response.Failed(err.Error(), c)
		return
	}
	response.SuccessNoData("文件权限修改成功", c)

}


