package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"

	"netdisk/global"
	"netdisk/model"
	"netdisk/model/requests"
	"netdisk/model/response"
	"netdisk/utils"
)


func checkShareAu(username,fileid,category string)bool {
	intFileid, _ := strconv.Atoi(fileid)
	if category == "1" {
		fv := new(model.FileVideo)
		fv.ID = uint(intFileid)
		fv.UploadUserName = username

		if !fv.IsExistFileVideo() || fv.Authority == model.High {
			return false
		}
		return true
	}

	fi := new(model.FileImage)
	fi.ID = uint(intFileid)
	fi.UploadUserName = username

	if !fi.IsExistFileImage() || fi.Authority == model.High {
		return false
	}
	return true
}


// @Summary 生成网盘分享二维码
// @Description  二维码中存放获取该分享文件信息链接
func (fs *fileService)Generateqrcode(Gq *requests.GQrcode,username string)(string,error) {
	if !checkShareAu(username,Gq.Fileid,Gq.Category){
		return "",errors.New("用户没有分享该文件权限")
	}

	admin := global.Config.Section("admin")
	domain := "http://" + admin.Key("host").String() + ":" + admin.Key("port").String()

	url := domain + "/share" + "/qrcode/?category=" + Gq.Category + "&fileid=" + Gq.Fileid + "&shareuser=" + username

	qrCode, err := qrcode.New(url, qrcode.Highest)
	if err != nil {
		global.SugaredLogger.Error("QR code generation failed", zap.Any("err:", err))
		return "", fmt.Errorf("二维码生成失败")
	}

	aliyun := utils.GetAliyunOss() //调用阿里云上传
	qrcodeUrl, err := aliyun.UploadQrcode(qrCode)
	if err != nil {
		global.SugaredLogger.Error("QR code cloud storage failed", zap.Any("err:", err))
		return "", fmt.Errorf("二维码生成失败")
	}
	return qrcodeUrl, nil
}


// @Summary 获取分享文件中的信息
// @Description  解析分享二维码中的链接
func (fs *fileService)GetFileInfoByQrcode(Aq *requests.AQrcode,username string)(*response.RespShareFileData,error) {
	intFileid, _ := strconv.Atoi(Aq.Fileid)
	switch Aq.Category {
	case "1":
		video := new(model.FileVideo)
		err := video.GetVideoInfo(intFileid)
		if err != nil {
			return nil, err
		}

		//权限验证
		if video.Authority==model.High{
			return nil,errors.New("该文件没有分享权限")
		}

		res := utils.FloatRound(float64(video.Size)/float64(1<<20), 3)
		shareFile := &response.RespShareFileData{
			FileId:     Aq.Fileid,
			FileName:   video.FileName,
			UploadTime: video.CreatedAt.Format("2006-01-02"),
			ShareUser:  Aq.Shareuser,
			Size:       strconv.FormatFloat(res, 'f', -1, 64) + "MB",
		}

		//用户获取分享文件后,将状态计入数据库
		vs := &model.VideoShare{
			Fileid:   uint(intFileid),
			Username: username,
		}

		if vs.CreateVideoShare() != nil {
			return nil, fmt.Errorf("保存分享信息失败")
		}
		return shareFile, nil

	case "2":
		image := new(model.FileImage)
		err := image.GetImageInfo(intFileid)
		if err != nil {
			return nil, err
		}

		if image.Authority==model.High{
			return nil,errors.New("该文件没有分享权限")
		}

		res := utils.FloatRound(float64(image.Size)/float64(1<<20), 3)
		shareFile := &response.RespShareFileData{
			FileId:     Aq.Fileid,
			FileName:   image.FileName,
			UploadTime: image.CreatedAt.Format("2006-01-02"),
			ShareUser:  Aq.Shareuser,
			Size:       strconv.FormatFloat(res, 'f', -1, 64) + "MB",
		}

		is := &model.ImageShare{
			Fileid:   uint(intFileid),
			Username: username,
		}

		if is.CreateImageShare() != nil {
			return nil, fmt.Errorf("保存分享信息失败")
		}
		return shareFile, nil
	default:
		return nil, errors.New("不支持的category值")
	}
}


// @Summary 生成加密链接
func (fs *fileService)GenerateSharingLink(gL *requests.GLink,username string)(string,error) {
	if !checkShareAu(username,gL.Fileid,gL.Category){
		return "",errors.New("用户没有分享该文件权限")
	}
	aL:= &requests.ALink{
		ShareUser:                username,
		Fileid:                   gL.Fileid,
		Category:                 gL.Category,
		EncryptionExtractionCode: gL.ExtractionCode,
	}
	//对提取码加密处理
	aL.EncryptionExtractionCode=utils.ExtractionCodeMD5Encryption(aL.EncryptionExtractionCode)

	//对信息结构体json序列化
	p,err:=json.Marshal(*aL)
	if err!=nil{
		global.SugaredLogger.Error("json marshal err",zap.Any("err",err))
		return "", fmt.Errorf("分享链接生成失败")
	}

	basep:=base64.URLEncoding.EncodeToString(p)

	admin := global.Config.Section("admin")
	domain := "http://" + admin.Key("host").String()
	url:=domain+"/s/"+basep
	return url, nil
}


// @Summary 获取分享文件中的信息
// @Description  解析分享链接
func (fs *fileService)GetFileInfoByLink(link,extractionCode,username string)(*response.RespShareFileData,error) {
	p,err:=base64.URLEncoding.DecodeString(link)
	if err!=nil{
		global.SugaredLogger.Error("base64.URLEncoding.DecodeString err:",zap.Any("err",err))
		return nil, fmt.Errorf("解析链接失败")
	}

	aL:=&requests.ALink{}
	err=json.Unmarshal(p,aL)
	if err!=nil{
		global.SugaredLogger.Error("json.Unmarshal err:",zap.Any("err",err))
		return nil, fmt.Errorf("解析链接失败")
	}

	//验证提取码正确性
	if aL.EncryptionExtractionCode!=utils.ExtractionCodeMD5Encryption(extractionCode){
		return nil, fmt.Errorf("提取码错误")
	}

	Aq:=&requests.AQrcode{
		Fileid: aL.Fileid,
		Category: aL.Category,
		Shareuser: aL.ShareUser,
	}

	return fs.GetFileInfoByQrcode(Aq,username)
}


