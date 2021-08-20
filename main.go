package main

import (
	"netdisk/core"
	"netdisk/global"
	inits "netdisk/initialization"
)

// @title Swagger Example API
// @version 1.0
// @description 接口文档

// @host xuzhihao.cn
// @BasePath
func main() {
	global.Config = inits.Ini()
	global.SugaredLogger = inits.Zap()
	global.Db = inits.MysqlGorm()
	global.Redis = inits.Redis()

	defer global.Db.DB().Close()

	if err := inits.InitTrans("zh"); err != nil {
		global.SugaredLogger.Panic(err)
	}
	core.RunServer()
}
