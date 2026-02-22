package ioc

import (
	"github.com/chiren-c/chili/pkg/loggerx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func InitLogger() loggerx.Logger {
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("logger", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	return loggerx.NewZapLogger(l)
}

func InitLoggerV1() loggerx.Logger {
	wd, _ := os.Getwd()
	lumberjackLogger := &lumberjack.Logger{
		// 注意有没有权限
		Filename:   wd + "/log/error.log", // 指定日志文件路径
		MaxSize:    50,                    // 每个日志文件的最大大小，单位：MB
		MaxBackups: 3,                     // 保留旧日志文件的最大个数
		MaxAge:     28,                    // 保留旧日志文件的最大天数
		Compress:   true,                  // 是否压缩旧的日志文件
	}

	// 创建zap日志核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		zapcore.DebugLevel, // 设置日志级别
	)

	l := zap.New(core, zap.AddCaller())
	return loggerx.NewZapLogger(l)
}
