package service

import (
	"errors"
	"strconv"
	"time"


	"netdisk/model"
	"netdisk/utils"
)

type userService struct{}

func GetUserService()*userService{return new(userService)}

func (u *userService)Register(username,password string)error{
	user:=&model.User{Username: username,Password: password}
	if user.IsExist(){
		return errors.New("该用户名已注册")
	}

	user.Salt = strconv.FormatInt(time.Now().Unix(), 10)       //通过时间戳获取撒盐值
	user.Password = utils.Md5Encryption(user.Password, user.Salt)    //md5加密


	return user.Create()
}


func (u *userService)Login(username,password string)(*model.User,error){
	user:=&model.User{Username: username,Password: password}
	if !user.IsExist(){
		return nil,errors.New("该用户名未注册")
	}

	password = utils.Md5Encryption(password, user.Salt)
	if password!=user.Password{
		return nil, errors.New("用户密码不正确")
	}
	return user,nil
}