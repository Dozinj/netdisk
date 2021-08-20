package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"mime/multipart"

	"go.uber.org/zap"

	"netdisk/global"
)

// @Summary md5用户密码加密
// @Description 撒盐值为服务层生成的时间戳
func PasswordMd5Encryption(password,salt string)string{
	m:=md5.New()
	m.Write([]byte(password))
	m.Write([]byte(salt))

	return hex.EncodeToString(m.Sum(nil))
}


// @Summary md5文件信息效验
// @Description  防止文件被篡改
func ContentMD5(file *multipart.FileHeader)string{
	m:=md5.New()
	fd,err:=file.Open()
	if err!=nil{
		global.SugaredLogger.Errorf("文件打开失败",zap.Any("err:",err))
		return ""
	}

	p,err:=ioutil.ReadAll(fd)
	if err!=nil{
		global.SugaredLogger.Errorf("文件读取失败",zap.Any("err:",err))
		return ""
	}
	m.Write(p)
	return base64.StdEncoding.EncodeToString(m.Sum(nil))
}

// @Summary 提取码加密
// @Description  使用配置中的jwtkey撒盐
func ExtractionCodeMD5Encryption(extractionCode string)string{
	m:=md5.New()
	if extractionCode!="" {
		m.Write([]byte(extractionCode))
	}
	m.Write([]byte(global.Config.Section("").Key("jwtkey").String()))

	return base64.StdEncoding.EncodeToString(m.Sum(nil))
}


