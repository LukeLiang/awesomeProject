package global

import (
	"awesomeProject/homework04/common"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	Logger *zap.Logger

	Config *common.Config
)
