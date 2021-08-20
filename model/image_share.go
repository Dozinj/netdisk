package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"netdisk/global"
)

type ImageShare struct {
	gorm.Model
	Fileid uint `json:"fileid"`
	Username string `json:"username"`
}

func (is *ImageShare)CreateImageShare()error{
	var is2 ImageShare
	if err:=global.Db.Where(is).First(&is2).Error;err== nil {
		//修改updateTime
		if err := global.Db.Model(&ImageShare{}).Where(is2.ID).Update("updated_at", time.Now()).Error; err != nil {
			global.SugaredLogger.Error("create video_share err:", zap.Any("err", err))
			return err
		}
	}else {
		//创建新纪录
		if err := global.Db.Create(is).Error; err != nil {
			global.SugaredLogger.Error("create image_share err:", zap.Any("err", err))
			return err
		}
	}
	return nil
}


func (is *ImageShare)IsExistImageShare()bool{
	if err:=global.Db.Where(is).First(is).Error;err!=nil{
		global.SugaredLogger.Error(err)
		return false
	}
	return true
}
