package initialization

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap()*zap.SugaredLogger{
	//WriterSyncer ：指定日志将写到哪里去。
	var lumberJackLogger = &lumberjack.Logger{
		Filename:   "./zap.log",
		MaxSize:    1,     //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,     //保留旧文件的最大个数
		MaxAge:     30,    //保留旧文件的最大天数
		Compress:   false, //是否压缩/归档旧文件
	}
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(lumberJackLogger)}
	//writes = append(writes, zapcore.AddSync(os.Stdout))  //关闭控制台输出

	//Encoder:编码器(如何写入日志)。
	encoderConfig:=zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder//修改时间编码器
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder//在日志文件中使用大写字母记录日志级别
	encoder:=zapcore.NewConsoleEncoder(encoderConfig)


	//func New(core zapcore.Core, options ...Option) *Logger
	//通过zapcore.NewMultiWriteSyncer(writes...)来设置多个输出
	core:=zapcore.NewCore(encoder,zapcore.NewMultiWriteSyncer(writes...),zapcore.DebugLevel)


	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	sugarLogger:=zap.New(core,caller).Sugar()
	return sugarLogger
}
