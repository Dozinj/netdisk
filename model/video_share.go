package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"netdisk/global"
)

type VideoShare struct {
	gorm.Model
	Fileid uint `json:"fileid"`
	Username string `json:"username"`
}

func (vs *VideoShare)CreateVideoShare()error {
	var vs2 VideoShare
	if global.Db.Where(vs).First(&vs2).Error == nil {
		//修改updateTime
		if err:=global.Db.Model(&VideoShare{}).Where(vs2.ID).Update("updated_at",time.Now()).Error;err!=nil{
			global.SugaredLogger.Error("create video_share err:", zap.Any("err", err))
			return err
		}
	} else {
		//创建一条新记录
		if err := global.Db.Create(vs).Error; err != nil {
			global.SugaredLogger.Error("create video_share err:", zap.Any("err", err))
			return err
		}
	}
	return nil
}


func (vs *VideoShare)IsExistVideoShare()bool {
	if err := global.Db.Where(vs).First(vs).Error; err != nil {
		global.SugaredLogger.Error(err)
		return false
	}
	return true
}


