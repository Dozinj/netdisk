package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"netdisk/global"
	"netdisk/model"
	"netdisk/model/requests"
	"netdisk/model/response"
	"netdisk/utils"
)
var err error


type fileService struct{}

func GetfileService()*fileService{return new(fileService)}

// @Summary 避免重传(秒传)
// @Description
func NoRetransmissionVideo(fileHeader *multipart.FileHeader)(string,string) {
	strMD5 := utils.ContentMD5(fileHeader) //加密文件信息
	//
	fileType := utils.GetFormat(fileHeader.Filename)

	redis := utils.GetRedis()
	if utils.CheckImageType(fileType) { //图片
		if redis.HExists("image", strMD5) {
			return redis.HGet("image", strMD5), strMD5
		} else {
			return "", strMD5
		}
	}

	//视频
	if redis.HExists("video", strMD5) {
		return redis.HGet("video", strMD5), strMD5
	} else {
		return "", strMD5
	}

}


// @Summary 上传视频和封面
// @Description
func(fs *fileService)UploadVideo(videoHeader *multipart.FileHeader,username string)(*model.FileVideo,error){
	var method uint
	if videoHeader.Size>50*(1<<20) {
		method = utils.NormalUpload
	} else {
		method = utils.AppendUpload
	}

	if videoHeader.Size > (1<<30) {
		return nil, fmt.Errorf("视频大于1GB")
	}
	videoFormat := utils.GetFormat(videoHeader.Filename)
	fmt.Println(videoHeader)
	if !utils.CheckVideoType(videoFormat){
		return nil, fmt.Errorf("不支持的视频格式")
	}


	//视频上传相关
	videoUrl,strMD5 := NoRetransmissionVideo(videoHeader)
	file := &model.FileVideo{
		FileName: videoHeader.Filename,
		FileUrl: videoUrl,
		Size:     videoHeader.Size,
		UploadUserName: username,
		FilePath: "/",
		Authority: model.Medium,
		MD5: strMD5,
	}
	if file.FileUrl!=""{
		//查询本人是否上传过相同的资源
		tmpfile :=&model.FileVideo{FileUrl: file.FileUrl,UploadUserName: username}
		if tmpfile.IsExistFileVideo(){
			return nil,fmt.Errorf("该视频已上传过")
		}
	}

	var url string
	if videoUrl == "" {
		//上传视频至oss
		aliyunOss:= utils.GetAliyunOss()
		switch method {
		case utils.NormalUpload:
			url, err = aliyunOss.UploadFile(videoHeader,strMD5)
		case utils.AppendUpload:
			url, err = aliyunOss.AppendUpload(videoHeader,strMD5)
		}

		if err!=nil {
			return nil, err
		}
		file.FileUrl=url
	}

	//该文件别人已上传情况,在upload_username中添加本人即可
	file.Model = gorm.Model{}
	err = file.CreateVideo()
	if err != nil {
		return nil, err
	}

	redis:=utils.GetRedis()
	//将文件MD5效验值存入redis中  若redis中未存在
	if !redis.HExists("video",strMD5){
		if err=redis.HSet("video",strMD5,file.FileUrl);err!=nil{
			global.SugaredLogger.Error("redis HSet video err",zap.Any("err",err))
		}
	}
	return file, nil
}

// @Summary 上传图片
// @Description
func (fs *fileService)UploadImage(fileHeader *multipart.FileHeader,username string)(*model.FileImage,error) {
	coverFormat := utils.GetFormat(fileHeader.Filename)
	if !utils.CheckImageType(coverFormat) {
		return nil, fmt.Errorf("不支持的图片格式")
	}

	aliyun := utils.GetAliyunOss() //调用阿里云上传
	imageUrl, strMD5 := NoRetransmissionVideo(fileHeader)

	image := &model.FileImage{
		FileName:       fileHeader.Filename,
		Size:           fileHeader.Size,
		FileUrl:        imageUrl,
		UploadUserName: username,
		FilePath: "/",
		Authority: model.Medium,
		MD5: strMD5,
	}
	if image.FileUrl != "" {
		tmpImage := &model.FileImage{FileUrl: image.FileUrl, UploadUserName: username}
		if tmpImage.IsExistFileImage() {
			return nil, fmt.Errorf("该图片已上传过")
		}
	}

	var url string
	if imageUrl == "" {
		//记录中查询不到则上传,获取该文件url
		url, err = aliyun.UploadFile(fileHeader, strMD5)
		if err != nil {
			return nil, err
		}
		image.FileUrl = url
	}

	//该文件别人已上传情况,在upload_username中添加本人即可
	image.Model = gorm.Model{}
	err = image.CreateImage()
	if err != nil {
		return nil, err
	}

	redis:=utils.GetRedis()
	if !redis.HExists("image",strMD5) {
		if err = redis.HSet("image",strMD5,image.FileUrl); err != nil {
			global.SugaredLogger.Error("redis HSet image err", zap.Any("err", err))
		}
	}
	return image, nil
}



// @Summary 保存网盘
// @Description  将username添加到该资源的upload_username中
func (fs *fileService)DiskSave(fc *requests.FC,username string)error {
	intFileid, _ := strconv.Atoi(fc.Fileid)
	switch fc.Category {
	case "1":
		fv := new(model.FileVideo)
		if err:=fv.GetVideoInfo(intFileid);err!=nil{
			return err
		}

		MD5str:=fv.MD5
		fv=&model.FileVideo{MD5: MD5str,UploadUserName: username}
		if fv.IsExistFileVideo() {
			return errors.New("该文件已在用户网盘中")
		}

		vs := &model.VideoShare{Username: username, Fileid: uint(intFileid)}

		if vs.IsExistVideoShare() { //是否有该文件保存权限(该文件是否已被分享给自己)
			fv := new(model.FileVideo)
			return fv.CreateVideoBySave(uint(intFileid), username)
		}

		return errors.New("用户没有该文件权限")
	case "2":
		fi := new(model.FileImage)
		if err:=fi.GetImageInfo(intFileid);err!=nil{
			return err
		}
		strMd5:=fi.MD5

		fi=&model.FileImage{MD5: strMd5,UploadUserName: username}
		if fi.IsExistFileImage() {
			return errors.New("该文件已在用户网盘中")
		}

		is := &model.ImageShare{Username: username, Fileid: uint(intFileid)}

		if is.IsExistImageShare() {
			fi := new(model.FileImage)
			return fi.CreateImageBySave(uint(intFileid), username)
		}

		return errors.New("用户没有该文件权限")
	}
	return errors.New("不支持的category值")
}


func (fs *fileService)LocalSave(fc *requests.FC,username string)(string,error) {
	intFileid, _ := strconv.Atoi(fc.Fileid)
	switch fc.Category {
	case "1":
		vs := &model.VideoShare{Username: username, Fileid: uint(intFileid)}
		ok1 := vs.IsExistVideoShare() //是否有该文件保存权限(该文件是否已被分享给自己)

		fv := &model.FileVideo{UploadUserName: username}
		fv.ID = uint(intFileid)
		ok2 := fv.IsExistFileVideo()
		if ok1 || ok2 {
			fv.UploadUserName = ""
			err := fv.GetVideoInfo(intFileid)
			if err != nil {
				return "", err
			}
			return fv.FileUrl, nil
		}
		return "",errors.New("用户没有该文件权限")

	case "2":
		is := &model.ImageShare{Username: username, Fileid: uint(intFileid)}
		//文件是否被分享给自己
		ok1 := is.IsExistImageShare()

		//或是文件是否在自己网盘中
		fi := &model.FileImage{UploadUserName: username}
		fi.ID = uint(intFileid)
		ok2 := fi.IsExistFileImage()

		if ok1 || ok2 {
			fi.UploadUserName = ""
			err := fi.GetImageInfo(intFileid)
			if err != nil {
				return "", err
			}
			return fi.FileUrl, nil
		}
		return "",errors.New("用户没有该文件权限")
	}
	return "", errors.New("不支持的category值")
}


func (fs *fileService)ListVideo(path,username string)(*response.RespListData,error) {
	//获取文件
	fv := &model.FileVideo{UploadUserName: username, FilePath: path}
	fvs, err := fv.ListAllVideo()
	if err != nil {
		return nil, err
	}

	files := make([]response.RespListFileData, 0, len(*fvs))
	for _, fv := range *fvs {
		res := utils.FloatRound(float64(fv.Size)/float64(1<<20), 3)
		data := response.RespListFileData{
			Fileid: fv.ID,
			CreatedAt: fv.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt: fv.UpdatedAt.Format("2006-01-02 15:04"),
			FileName:  fv.FileName,
			Size:      strconv.FormatFloat(res, 'f', -1, 64) + "MB",
			Authority: fv.Authority,
		}
		files = append(files, data)
	}

	//获取文件夹
	vd := &model.VideoFolder{FolderPath: path, Username: username}
	vds, err := vd.ListFolders()
	if err != nil {
		return nil, err
	}
	folders := make([]response.RespListFolderData, 0, len(*vds))
	for _, vd := range *vds {
		data := response.RespListFolderData{
			CreatedAt:  vd.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt:  vd.UpdatedAt.Format("2006-01-02 15:04"),
			FolderName: vd.FolderName,
		}
		folders = append(folders, data)
	}
	ListData := &response.RespListData{File: files, Folder: folders}
	return ListData, nil
}


func (fs *fileService)ListImage(path,username string)(*response.RespListData,error) {
	//获取文件
	fi := &model.FileImage{UploadUserName: username, FilePath: path}
	fis, err := fi.ListAllImage()
	if err != nil {
		return nil, err
	}

	files := make([]response.RespListFileData, 0, len(*fis))
	for _, fi := range *fis {
		res := utils.FloatRound(float64(fi.Size)/float64(1<<20), 3)
		data := response.RespListFileData{
			Fileid: fi.ID,
			CreatedAt: fi.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt: fi.UpdatedAt.Format("2006-01-02 15:04"),
			FileName:  fi.FileName,
			Size:      strconv.FormatFloat(res, 'f', -1, 64) + "MB",
			Authority: fi.Authority,
		}
		files = append(files, data)
	}

	//获取文件夹
	iF := &model.ImageFolder{FolderPath: path, Username: username}
	iFs, err := iF.ListFolders()
	if err != nil {
		return nil, err
	}
	folders := make([]response.RespListFolderData, 0, len(*iFs))
	for _, iF := range *iFs {
		data := response.RespListFolderData{
			CreatedAt:  iF.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt:  iF.UpdatedAt.Format("2006-01-02 15:04"),
			FolderName: iF.FolderName,
		}
		folders = append(folders, data)
	}
	ListData := &response.RespListData{File: files, Folder: folders}
	return ListData, nil
}



func (fs *fileService)ChangeVideoPath(fileid,filePath,username string)error {
	//验证用户是否有该文件权限
	fv := new(model.FileVideo)
	intFileid, _ := strconv.Atoi(fileid)
	fv.ID = uint(intFileid)
	fv.UploadUserName = username

	if !fv.IsExistFileVideo() {
		return errors.New("该用户没有此文件权限")
	}

	if err := createVPath(filePath, username); err != nil {
		return err
	}

	return fv.UpdateFilePath(filePath)
}


func createVPath(filePath,username string)error {
	filePaths := strings.Split(filePath, "/")

	var path string
	//获取多级目录
	for _, tmppath := range filePaths {
		if tmppath == "" {
			path = "/"
		} else {
			path = path + tmppath + "/"
		}

		//判断该层目录是否存在
		fv := new(model.FileVideo)
		fv.UploadUserName = username
		fv.FilePath = path

		ok1 := fv.IsExistFileVideo()

		vf := new(model.VideoFolder)
		vf.Username = username
		vf.FolderPath = path

		ok2 := vf.IsExist()

		if !ok1 && !ok2 {
			//新建文件夹
			folderPath := strings.Split(path, "/"+tmppath)[0]
			vf.FolderPath=folderPath+"/"
			vf.FolderName = tmppath

			if err := vf.Create(); err != nil {
				return err
			}
		}
	}
	return nil
}


func (fs *fileService)ChangeImagePath(fileid,filePath,username string)error {
	//验证用户是否有该文件权限
	fi := new(model.FileImage)
	intFileid, _ := strconv.Atoi(fileid)
	fi.ID = uint(intFileid)
	fi.UploadUserName = username

	if !fi.IsExistFileImage() {
		return errors.New("该用户没有此文件权限")
	}

	if err := createIPath(filePath, username); err != nil {
		return err
	}

	return fi.UpdateFilePath(filePath)
}


func createIPath(filePath,username string)error {
	filePaths := strings.Split(filePath, "/")

	var path string
	//获取多级目录
	for _, tmppath := range filePaths {
		if tmppath == "" {
			path = "/"
		} else {
			path = path + tmppath + "/"
		}

		//判断该层目录是否存在
		fi := new(model.FileImage)
		fi.UploadUserName = username
		fi.FilePath = path

		ok1 := fi.IsExistFileImage()

		If := new(model.ImageFolder)
		If.Username = username
		If.FolderPath = path

		ok2 := If.IsExist()

		if !ok1 && !ok2 {
			//新建文件夹
			folderPath := strings.Split(path, "/"+tmppath)[0]
			If.FolderPath=folderPath+"/"
			If.FolderName = tmppath

			if err := If.Create(); err != nil {
				return err
			}
		}
	}
	return nil
}


func (fs *fileService)ChangeFile(C interface{},fileid string,category string,username string)error {
	switch category {
	case "1":
		fv := new(model.FileVideo)
		fv.UploadUserName = username
		intFileid, _ := strconv.Atoi(fileid)
		fv.ID = uint(intFileid)

		if !fv.IsExistFileVideo() {
			return errors.New("该用户没有操作此文件权限")
		}


		switch reflect.TypeOf(C).Name() {
		case "ChangeFilename":
			Cf := C.(requests.ChangeFilename)
			return changeVideoName(fv, Cf.Filename)
		case "ChangeAuthority":
			Ca := C.(requests.ChangeAuthority)
			return changeVideoAu(fv, Ca.Authority)
		}

	case "2":
		fi := new(model.FileImage)
		fi.UploadUserName = username
		intFileid, _ := strconv.Atoi(fileid)
		fi.ID = uint(intFileid)

		if !fi.IsExistFileImage() {
			return errors.New("该用户没有操作此文件权限")
		}

		switch reflect.TypeOf(C).Name() {
		case "ChangeFilename":
			Cf := C.(requests.ChangeFilename)
			return changeImageName(fi, Cf.Filename)
		case "ChangeAuthority":
			Ca := C.(requests.ChangeAuthority)
			return changeImageAu(fi, Ca.Authority)
		}
	}

	return errors.New("not support category")
}

func changeVideoName(f *model.FileVideo,filename string)error{
	filenames := strings.Split(f.FileName, ".")
	format := filenames[len(filenames)-1]

	newFilename := filename + "." + format
	return f.UpdateFileName(newFilename)
}

func changeImageName(f *model.FileImage,filename string)error{
	filenames := strings.Split(f.FileName, ".")
	format := filenames[len(filenames)-1]

	newFilename := filename + "." + format
	return f.UpdateFileName(newFilename)
}

func changeVideoAu(f *model.FileVideo,authority string)error{
	intAuthority,err:=strconv.Atoi(authority)
	if err!=nil{
		global.SugaredLogger.Error("类型转换失败",zap.Any("err:",err))
		return errors.New("修改失败")
	}
	return f.UpdateAuthority(intAuthority)
}

func changeImageAu(f *model.FileImage,authority string)error{
	intAuthority,err:=strconv.Atoi(authority)
	if err!=nil{
		global.SugaredLogger.Error("类型转换失败",zap.Any("err:",err))
		return errors.New("修改失败")
	}
	return f.UpdateAuthority(intAuthority)
}


