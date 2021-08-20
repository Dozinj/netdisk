package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"netdisk/global"
)

type User struct {
	gorm.Model
	Username  string `json:"username" example:"xzh"`     // 用户登录名
	Password  string `json:"password"  example:"123456"` // 用户登录密码
	Salt      string `json:"salt"`                       //用户密码加密撒盐值
}


func (u *User)IsExist()bool{
	if !errors.Is(global.Db.Where("username=?", u.Username).First(u).Error, gorm.ErrRecordNotFound){
		return true
	}
	return false
}

func (u *User)Create()error{
	if err:= global.Db.Create(u).Error;err!=nil{
		global.SugaredLogger.Error("用户信息储存失败",zap.Any("err:",err))
		return errors.New("用户信息储存失败")
	}
	return nil
}

