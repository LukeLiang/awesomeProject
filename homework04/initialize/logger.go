package initialize

import (
	"awesomeProject/homework04/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式化
	var err error
	var logger *zap.Logger
	logger, err = config.Build()
	if err != nil {
		panic(err)
	}

	global.Logger = logger
}
