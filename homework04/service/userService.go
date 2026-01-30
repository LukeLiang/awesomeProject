package service

import (
	"awesomeProject/homework04/global"
	"awesomeProject/homework04/model"

	"go.uber.org/zap"
)

type UserService struct {
}

// 确保 UserService 实现 IUserService 接口
var _ IUserService = (*UserService)(nil)

func (s *UserService) Create(user *model.User) error {
	global.Logger.Debug("Creating user",
		zap.String("username", user.Username),
	)

	if err := global.DB.Create(&user).Error; err != nil {
		global.Logger.Error("Failed to create user",
			zap.String("username", user.Username),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("User created successfully",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)
	return nil
}

func (s *UserService) FindByUsername(username string) (*model.User, error) {
	global.Logger.Debug("Finding user by username",
		zap.String("username", username),
	)

	var user model.User
	result := global.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		global.Logger.Warn("User not found",
			zap.String("username", username),
			zap.Error(result.Error),
		)
		return nil, result.Error
	}

	global.Logger.Debug("User found",
		zap.Uint("user_id", user.ID),
		zap.String("username", username),
	)
	return &user, nil
}

func (s *UserService) FindByID(id uint) (*model.User, error) {
	global.Logger.Debug("Finding user by ID",
		zap.Uint("user_id", id),
	)

	var user model.User
	result := global.DB.First(&user, id)
	if result.Error != nil {
		global.Logger.Warn("User not found",
			zap.Uint("user_id", id),
			zap.Error(result.Error),
		)
		return nil, result.Error
	}

	global.Logger.Debug("User found",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)
	return &user, nil
}
