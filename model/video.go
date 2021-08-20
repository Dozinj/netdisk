package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"netdisk/global"
)

const (
	Low = 8+iota //所有人可下载
	Medium		//仅获得链接的人可下载
	High        //仅自己见
)

//视频相关
type FileVideo struct {
	gorm.Model
	FileName       string `json:"file_name"`        //视频名称
	FileUrl        string `json:"video_url"`        //视频资源地址
	Size           int64  `json:"size"`             //视频大小
	UploadUserName string `json:"upload_user_name"` //上传者姓名
	MD5            string `json:"md_5"`             //文件内容MD5加密
	FilePath       string `json:"file_path"`        //网盘路径 (默认根目录)
	Authority      int    `json:"authority"`        //文件权限(所有人可下载、仅获得链接的人可下载、仅自己见，默认medium)
}

func (f *FileVideo)GetVideoInfo(fileid int)error {
	f.ID = uint(fileid)
	err := global.Db.Where(f).First(f).Error
	if err != nil {
		global.SugaredLogger.Error("failed to query file", zap.Any("err:", err))
		return fmt.Errorf("未找到该文件信息")
	}
	return nil
}

func (f *FileVideo)CreateVideo()error{
	if err := global.Db.Create(f).Error; err != nil {
		global.SugaredLogger.Error(err)
		return  fmt.Errorf("储存信息失败")
	}
	return nil
}

func (f *FileVideo)CreateVideoBySave(fileid uint,username string)error{
	//查询出原文件信息
	f.ID=fileid
	if err:=global.Db.Where(f.ID).First(f).Error;err!=nil{
		global.SugaredLogger.Error("find video failed",zap.Any("err",err))
		return fmt.Errorf("保存到网盘失败")
	}

	fv2:=f
	fv2.UploadUserName=username
	fv2.Model=gorm.Model{}
	return fv2.CreateVideo()
}


func (f *FileVideo)UpdateFilePath(filepath string)error {
	if err := global.Db.Model(&FileVideo{}).Where(f).Update("file_path",filepath).Error; err != nil {
		global.SugaredLogger.Error("update failed", zap.Any("err", err))
		return fmt.Errorf("更新到网盘失败")
	}
	return nil
}


func (f *FileVideo)UpdateFileName(filename string)error {
	if err := global.Db.Model(&FileVideo{}).Where(f).Update("file_name",filename).Error; err != nil {
		global.SugaredLogger.Error("update failed", zap.Any("err", err))
		return fmt.Errorf("更新到网盘失败")
	}
	return nil
}

func (f *FileVideo)UpdateAuthority(authority int)error {
	if err := global.Db.Model(&FileVideo{}).Where(f).Update("authority",authority).Error; err != nil {
		global.SugaredLogger.Error("update failed", zap.Any("err", err))
		return fmt.Errorf("更新到网盘失败")
	}
	return nil
}


func (f *FileVideo)IsExistFileVideo()bool{
	if err:=global.Db.Where(f).First(f).Error;err!=nil{
		global.SugaredLogger.Info(err)
		return false
	}
	return true
}

func (f *FileVideo)ListAllVideo()(*[]FileVideo,error ){
	fvs:=new([]FileVideo)
	if err:=global.Db.Where(f).Find(fvs).Error;err!=nil{
		global.SugaredLogger.Info(err)
		return nil,fmt.Errorf("查询信息失败")
	}
	return fvs,nil
}
