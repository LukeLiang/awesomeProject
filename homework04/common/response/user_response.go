package response

import "awesomeProject/homework04/common/types"

// UserResponse 用户响应（不含密码）
type UserResponse struct {
	ID        uint            `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	CreatedAt types.LocalTime `json:"createdAt"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user,omitempty"`
}
