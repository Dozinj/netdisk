package initialization

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"netdisk/global"
	"netdisk/model"
)

func MysqlGorm()*gorm.DB {
	m := global.Config.Section("mysql")
	db, err := gorm.Open(m.Key("dialect").String(),
		m.Key("user").String()+":"+m.Key("password").String()+
			"@tcp("+m.Key("host").String()+":"+m.Key("port").String()+")/"+
			m.Key("database").String()+"?charset=utf8mb4&parseTime=True&loc=Local")

	if err != nil {
		global.SugaredLogger.Fatal("数据库驱动失败", err)
		return nil
	}

	db.LogMode(true) //打印日志
	CreateTable(db)
	return db
}

func CreateTable(db *gorm.DB){
	err:=db.AutoMigrate(
		&model.User{},
		&model.FileVideo{},
		&model.FileImage{},
		&model.ImageShare{},
		&model.VideoShare{},
		&model.VideoFolder{},
		&model.ImageFolder{},
	).Error

	if err!=nil{
		global.SugaredLogger.Panic("创建数据表失败",zap.Any("err:",err))
		return
	}
	global.SugaredLogger.Info("创建数据表成功")
}
