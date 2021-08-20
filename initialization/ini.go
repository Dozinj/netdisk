package initialization

import (
	"gopkg.in/ini.v1"
)

func Ini()*ini.File {
	config, err := ini.Load("./conf/config.ini")

	if err != nil {
		panic(err)
	}
	return config
}
