package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gopkg.in/ini.v1"
)

var (
	Db            *gorm.DB
	Redis         *redis.Client
	Config        *ini.File
	SugaredLogger *zap.SugaredLogger
	Trans         ut.Translator
)
