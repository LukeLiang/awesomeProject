package service

import (
	"awesomeProject/homework04/common/request"
	"awesomeProject/homework04/model"
)

// IUserService 用户服务接口
type IUserService interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
}

// IPostService 文章服务接口
type IPostService interface {
	Create(post *model.Posts) error
	Update(post *model.Posts) error
	Delete(id uint) error
	GetByID(id uint) (*model.Posts, error)
	List(userId uint, page *request.PageRequest) ([]model.Posts, int64, error)
}

// ICommentService 评论服务接口
type ICommentService interface {
	Create(comment *model.Comments) error
	Delete(id uint) error
	GetByID(id uint) (*model.Comments, error)
	ListByPostID(postId uint, page *request.PageRequest) ([]model.Comments, int64, error)
}
