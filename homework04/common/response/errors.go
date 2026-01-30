package response

import "net/http"

// AppError 应用错误结构
type AppError struct {
	Code       int    // 业务错误码
	Message    string // 错误消息
	HttpStatus int    // HTTP 状态码
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError 创建新的应用错误
func NewAppError(code int, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HttpStatus: httpStatus,
	}
}

// 通用错误 (1000-1999)
var (
	ErrInvalidParams  = &AppError{Code: 1001, Message: "参数错误", HttpStatus: http.StatusBadRequest}
	ErrServerInternal = &AppError{Code: 1002, Message: "服务器内部错误", HttpStatus: http.StatusInternalServerError}
	ErrNotFound       = &AppError{Code: 1003, Message: "资源不存在", HttpStatus: http.StatusNotFound}
	ErrMethodNotAllow = &AppError{Code: 1004, Message: "请求方法不允许", HttpStatus: http.StatusMethodNotAllowed}
)

// 认证错误 (2000-2999)
var (
	ErrUnauthorized     = &AppError{Code: 2000, Message: "未授权", HttpStatus: http.StatusUnauthorized}
	ErrTokenInvalid     = &AppError{Code: 2001, Message: "Token 无效或已过期", HttpStatus: http.StatusUnauthorized}
	ErrTokenFormat      = &AppError{Code: 2002, Message: "Token 格式错误", HttpStatus: http.StatusUnauthorized}
	ErrPermissionDenied = &AppError{Code: 2003, Message: "权限不足", HttpStatus: http.StatusForbidden}
)

// 用户错误 (3000-3999)
var (
	ErrUserNotFound       = &AppError{Code: 3000, Message: "用户不存在", HttpStatus: http.StatusNotFound}
	ErrUserExists         = &AppError{Code: 3001, Message: "用户已存在", HttpStatus: http.StatusBadRequest}
	ErrUsernameOrPassword = &AppError{Code: 3002, Message: "用户名或密码错误", HttpStatus: http.StatusBadRequest}
	ErrPasswordHash       = &AppError{Code: 3003, Message: "密码加密失败", HttpStatus: http.StatusInternalServerError}
	ErrUserCreate         = &AppError{Code: 3004, Message: "用户创建失败", HttpStatus: http.StatusInternalServerError}
)

// 文章错误 (4000-4999)
var (
	ErrPostNotFound = &AppError{Code: 4000, Message: "文章不存在", HttpStatus: http.StatusNotFound}
	ErrPostCreate   = &AppError{Code: 4001, Message: "文章创建失败", HttpStatus: http.StatusInternalServerError}
	ErrPostUpdate   = &AppError{Code: 4002, Message: "文章更新失败", HttpStatus: http.StatusInternalServerError}
	ErrPostDelete   = &AppError{Code: 4003, Message: "文章删除失败", HttpStatus: http.StatusInternalServerError}
)

// 评论错误 (5000-5999)
var (
	ErrCommentNotFound = &AppError{Code: 5000, Message: "评论不存在", HttpStatus: http.StatusNotFound}
	ErrCommentCreate   = &AppError{Code: 5001, Message: "评论创建失败", HttpStatus: http.StatusInternalServerError}
)

// 数据库错误 (6000-6999)
var (
	ErrDBConnection = &AppError{Code: 6000, Message: "数据库连接失败", HttpStatus: http.StatusInternalServerError}
	ErrDBQuery      = &AppError{Code: 6001, Message: "数据库查询失败", HttpStatus: http.StatusInternalServerError}
	ErrDBCreate     = &AppError{Code: 6002, Message: "数据库创建失败", HttpStatus: http.StatusInternalServerError}
	ErrDBUpdate     = &AppError{Code: 6003, Message: "数据库更新失败", HttpStatus: http.StatusInternalServerError}
	ErrDBDelete     = &AppError{Code: 6004, Message: "数据库删除失败", HttpStatus: http.StatusInternalServerError}
)
