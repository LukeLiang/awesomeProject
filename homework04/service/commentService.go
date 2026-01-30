package service

import (
	"awesomeProject/homework04/common/request"
	"awesomeProject/homework04/global"
	"awesomeProject/homework04/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentService struct {
}

// 确保 CommentService 实现 ICommentService 接口
var _ ICommentService = (*CommentService)(nil)

func (s *CommentService) Create(comment *model.Comments) error {
	global.Logger.Debug("Creating comment",
		zap.Uint("post_id", comment.PostId),
		zap.Uint("user_id", comment.UserId),
	)

	if err := global.DB.Create(comment).Error; err != nil {
		global.Logger.Error("Failed to create comment",
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("Comments created successfully",
		zap.Uint("comment_id", comment.ID),
	)
	return nil
}

func (s *CommentService) Delete(id uint) error {
	global.Logger.Debug("Deleting comment",
		zap.Uint("comment_id", id),
	)

	result := global.DB.Delete(&model.Comments{}, id)
	if result.Error != nil {
		global.Logger.Error("Failed to delete comment",
			zap.Uint("comment_id", id),
			zap.Error(result.Error),
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		global.Logger.Warn("Comments not found for delete",
			zap.Uint("comment_id", id),
		)
		return gorm.ErrRecordNotFound
	}

	global.Logger.Info("Comments deleted successfully",
		zap.Uint("comment_id", id),
	)
	return nil
}

func (s *CommentService) GetByID(id uint) (*model.Comments, error) {
	global.Logger.Debug("Getting comment by ID",
		zap.Uint("comment_id", id),
	)

	var comment model.Comments
	if err := global.DB.First(&comment, id).Error; err != nil {
		global.Logger.Warn("Comments not found",
			zap.Uint("comment_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return &comment, nil
}

func (s *CommentService) ListByPostID(postId uint, page *request.PageRequest) ([]model.Comments, int64, error) {
	global.Logger.Debug("Listing comments for post",
		zap.Uint("post_id", postId),
		zap.Int("page", page.GetPage()),
		zap.Int("pageSize", page.GetPageSize()),
	)

	var comments []model.Comments
	var total int64

	query := global.DB.Model(&model.Comments{}).Where("post_id = ?", postId)

	if err := query.Count(&total).Error; err != nil {
		global.Logger.Error("Failed to count comments",
			zap.Uint("post_id", postId),
			zap.Error(err),
		)
		return nil, 0, err
	}

	if err := query.Offset(page.GetOffset()).Limit(page.GetPageSize()).
		Order("created_at DESC").Find(&comments).Error; err != nil {
		global.Logger.Error("Failed to list comments",
			zap.Uint("post_id", postId),
			zap.Error(err),
		)
		return nil, 0, err
	}

	global.Logger.Debug("Comments retrieved",
		zap.Uint("post_id", postId),
		zap.Int64("total", total),
		zap.Int("count", len(comments)),
	)
	return comments, total, nil
}
