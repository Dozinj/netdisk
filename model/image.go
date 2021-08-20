package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"netdisk/global"
)

//图片相关
type FileImage struct {
	gorm.Model
	FileName       string `json:"file_name"`        //图片名称
	FileUrl        string `json:"image_url"`        //图片存储地址
	Size           int64  `json:"size"`             //图片内存大小
	UploadUserName string `json:"upload_user_name"` //上传者姓名
	MD5            string `json:"md_5"`             //文件内容MD5加密
	FilePath       string `json:"file_path"`        //网盘路径 (默认根目录)
	Authority      int    `json:"authority"`        //文件权限(所有人可下载、仅获得链接的人可下载、仅自己见，默认medium)
}

func (f *FileImage)GetImageInfo(fileid int)error {
	f.ID = uint(fileid)
	err := global.Db.Where(f).First(f).Error
	if err != nil {
		global.SugaredLogger.Error("Failed to query file", zap.Any("err:", err))
		return fmt.Errorf("未找到该文件信息")
	}
	return nil
}


func (f *FileImage)CreateImage()error {
	if err := global.Db.Create(f).Error; err != nil {
		global.SugaredLogger.Error(err)
		return fmt.Errorf("储存信息失败")
	}
	return nil
}


func (f *FileImage)CreateImageBySave(fileid uint,username string)error{
	//查询出原文件信息
	f.ID=fileid
	if err:=global.Db.Where(f.ID).First(f).Error;err!=nil{
		global.SugaredLogger.Error("find file image failed:",zap.Any("err",err))
		return fmt.Errorf("保存到网盘失败")
	}

	fi2:=f
	fi2.UploadUserName=username
	fi2.Authority=Medium
	fi2.FilePath="/"
	fi2.Model=gorm.Model{}  //保留关键信息,插入新记录
	return fi2.CreateImage()
}


func (f *FileImage)UpdateFilePath(filepath string)error {
	if err := global.Db.Model(&FileImage{}).Where(f).UpdateColumn("file_path",filepath).Error; err != nil {
		global.SugaredLogger.Error("update failed", zap.Any("err", err))
		return fmt.Errorf("更新到网盘失败")
	}
	return nil
}

func (f *FileImage)UpdateFileName(filename string)error {
	if err := global.Db.Model(&FileImage{}).Where(f).Update("file_name",filename).Error; err != nil {
		global.SugaredLogger.Error("update failed", zap.Any("err", err))
		return fmt.Errorf("更新到网盘失败")
	}
	return nil
}

func (f *FileImage)UpdateAuthority(authority int)error {
	if err := global.Db.Model(&FileImage{}).Where(f).Update("authority",authority).Error; err != nil {
		global.SugaredLogger.Error("update failed", zap.Any("err", err))
		return fmt.Errorf("更新到网盘失败")
	}
	return nil
}


func (f *FileImage)IsExistFileImage()bool{
	if err:=global.Db.Where(f).First(f).Error;err!=nil{
		global.SugaredLogger.Info(err)
		return false
	}
	return true
}


func (f *FileImage)ListAllImage()(*[]FileImage,error ){
	ivs:=new([]FileImage)
	if err:=global.Db.Where(f).Find(ivs).Error;err!=nil{
		global.SugaredLogger.Info(err)
		return nil,fmt.Errorf("查询信息失败")
	}
	return ivs,nil
}