package model

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"netdisk/global"
)

type VideoFolder struct {
	gorm.Model
	FolderPath string `json:"folder_path"` //文件夹路径
	FolderName string `json:"folder_name"` //文件夹名称
	Username   string `json:"username"`    //用户名
}

func (f *VideoFolder)IsExist()bool{
	if err:=global.Db.Where(f).First(f).Error;err!=nil{
		global.SugaredLogger.Info(err)
		return false
	}
	return true
}

func (f *VideoFolder)Create()error{
	if err := global.Db.Create(f).Error; err != nil {
		global.SugaredLogger.Error(err)
		return  fmt.Errorf("储存信息失败")
	}
	return nil
}

func (f *VideoFolder)ListFolders()(*[]VideoFolder,error){
	var vds []VideoFolder
	if err := global.Db.Where(f).Find(&vds).Error; err != nil {
		global.SugaredLogger.Error(err)
		return  nil,fmt.Errorf("储存信息失败")
	}
	return &vds,nil
}



