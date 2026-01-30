package api

import (
	"time"

	"awesomeProject/homework04/common/request"
	"awesomeProject/homework04/common/response"
	"awesomeProject/homework04/global"
	"awesomeProject/homework04/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserApi struct {
}

func (a *UserApi) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Warn("Register: invalid params",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		response.FailWithMessage(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 检查用户是否已存在
	existingUser, _ := userService.FindByUsername(req.Username)
	if existingUser != nil {
		global.Logger.Warn("Register: user already exists",
			zap.String("username", req.Username),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Fail(c, response.ErrUserExists)
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		global.Logger.Error("Register: password hash failed",
			zap.Error(err),
			zap.String("username", req.Username),
		)
		response.Fail(c, response.ErrPasswordHash)
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := userService.Create(user); err != nil {
		global.Logger.Error("Register: create user failed",
			zap.Error(err),
			zap.String("username", req.Username),
		)
		response.Fail(c, response.ErrUserCreate)
		return
	}

	global.Logger.Info("User registered",
		zap.String("username", user.Username),
		zap.Uint("user_id", user.ID),
	)
	response.SuccessWithData(c, "用户注册成功", response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (a *UserApi) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Warn("Login: invalid params",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		response.FailWithMessage(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 校验用户名和密码
	user, err := userService.FindByUsername(req.Username)
	if err != nil {
		global.Logger.Warn("Login: user not found",
			zap.String("username", req.Username),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Fail(c, response.ErrUsernameOrPassword)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		global.Logger.Warn("Login: password mismatch",
			zap.String("username", req.Username),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Fail(c, response.ErrUsernameOrPassword)
		return
	}

	// 登录成功 颁发JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(global.Config.Secret)
	if err != nil {
		global.Logger.Error("Login: JWT sign failed",
			zap.Error(err),
			zap.String("username", req.Username),
		)
		response.Fail(c, response.ErrServerInternal)
		return
	}

	global.Logger.Info("User logged in",
		zap.String("username", user.Username),
		zap.Uint("user_id", user.ID),
	)
	response.Success(c, response.LoginResponse{
		Token: tokenString,
		User: &response.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}
