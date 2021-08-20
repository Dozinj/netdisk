package core

import (
	"net/http"

	"go.uber.org/zap"

	"netdisk/global"
	"netdisk/initialization"
)

func RunServer(){
	router:= initialization.Routers()

	admin:=global.Config.Section("admin")
	address:=admin.Key("host").String()+":"+admin.Key("port").String()


	server:=&http.Server{
		Addr:           address,
		Handler:        router,
		//ReadTimeout:    30 * time.Second,
		//WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	global.SugaredLogger.Info("server run success on ", zap.String("address", address))
	global.SugaredLogger.Error(server.ListenAndServe().Error())
}
