package service

import (
	"awesomeProject/homework04/common/request"
	"awesomeProject/homework04/global"
	"awesomeProject/homework04/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostService struct {
}

// 确保 PostService 实现 IPostService 接口
var _ IPostService = (*PostService)(nil)

func (s *PostService) Create(post *model.Posts) error {
	global.Logger.Debug("Creating post",
		zap.String("title", post.Title),
		zap.Uint("user_id", post.UserId),
	)

	if err := global.DB.Create(post).Error; err != nil {
		global.Logger.Error("Failed to create post",
			zap.String("title", post.Title),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("Posts created successfully",
		zap.Uint("post_id", post.ID),
		zap.String("title", post.Title),
	)
	return nil
}

func (s *PostService) Update(post *model.Posts) error {
	global.Logger.Debug("Updating post",
		zap.Uint("post_id", post.ID),
	)

	result := global.DB.Model(post).Updates(map[string]interface{}{
		"title":   post.Title,
		"content": post.Content,
	})

	if result.Error != nil {
		global.Logger.Error("Failed to update post",
			zap.Uint("post_id", post.ID),
			zap.Error(result.Error),
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		global.Logger.Warn("Posts not found for update",
			zap.Uint("post_id", post.ID),
		)
		return gorm.ErrRecordNotFound
	}

	global.Logger.Info("Posts updated successfully",
		zap.Uint("post_id", post.ID),
	)
	return nil
}

func (s *PostService) Delete(id uint) error {
	global.Logger.Debug("Deleting post",
		zap.Uint("post_id", id),
	)

	result := global.DB.Delete(&model.Posts{}, id)
	if result.Error != nil {
		global.Logger.Error("Failed to delete post",
			zap.Uint("post_id", id),
			zap.Error(result.Error),
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		global.Logger.Warn("Posts not found for delete",
			zap.Uint("post_id", id),
		)
		return gorm.ErrRecordNotFound
	}

	global.Logger.Info("Posts deleted successfully",
		zap.Uint("post_id", id),
	)
	return nil
}

func (s *PostService) GetByID(id uint) (*model.Posts, error) {
	global.Logger.Debug("Getting post by ID",
		zap.Uint("post_id", id),
	)

	var post model.Posts
	if err := global.DB.First(&post, id).Error; err != nil {
		global.Logger.Warn("Posts not found",
			zap.Uint("post_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Debug("Posts found",
		zap.Uint("post_id", id),
		zap.String("title", post.Title),
	)
	return &post, nil
}

func (s *PostService) List(userId uint, page *request.PageRequest) ([]model.Posts, int64, error) {
	global.Logger.Debug("Listing posts for user",
		zap.Uint("user_id", userId),
		zap.Int("page", page.GetPage()),
		zap.Int("pageSize", page.GetPageSize()),
	)

	var posts []model.Posts
	var total int64

	query := global.DB.Model(&model.Posts{}).Where("user_id = ?", userId)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Error("Failed to count posts",
			zap.Uint("user_id", userId),
			zap.Error(err),
		)
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(page.GetOffset()).Limit(page.GetPageSize()).
		Order("created_at DESC").Find(&posts).Error; err != nil {
		global.Logger.Error("Failed to list posts",
			zap.Uint("user_id", userId),
			zap.Error(err),
		)
		return nil, 0, err
	}

	global.Logger.Debug("Posts retrieved",
		zap.Uint("user_id", userId),
		zap.Int64("total", total),
		zap.Int("count", len(posts)),
	)
	return posts, total, nil
}

// DeleteWithComments 删除文章及其所有评论（事务）
func (s *PostService) DeleteWithComments(id uint) error {
	global.Logger.Debug("Deleting post with comments",
		zap.Uint("post_id", id),
	)

	return global.Transaction(func(tx *gorm.DB) error {
		// 先删除所有评论
		if err := tx.Where("post_id = ?", id).Delete(&model.Comments{}).Error; err != nil {
			global.Logger.Error("Failed to delete comments",
				zap.Uint("post_id", id),
				zap.Error(err),
			)
			return err
		}

		// 再删除文章
		result := tx.Delete(&model.Posts{}, id)
		if result.Error != nil {
			global.Logger.Error("Failed to delete post",
				zap.Uint("post_id", id),
				zap.Error(result.Error),
			)
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		global.Logger.Info("Posts and comments deleted successfully",
			zap.Uint("post_id", id),
		)
		return nil
	})
}

// CreateWithComments 创建文章并批量创建评论（事务）
func (s *PostService) CreateWithComments(post *model.Posts, comments []model.Comments) error {
	global.Logger.Debug("Creating post with comments",
		zap.String("title", post.Title),
		zap.Int("comments_count", len(comments)),
	)

	return global.Transaction(func(tx *gorm.DB) error {
		// 创建文章
		if err := tx.Create(post).Error; err != nil {
			global.Logger.Error("Failed to create post in transaction",
				zap.Error(err),
			)
			return err
		}

		// 批量创建评论
		if len(comments) > 0 {
			for i := range comments {
				comments[i].PostId = post.ID
			}
			if err := tx.Create(&comments).Error; err != nil {
				global.Logger.Error("Failed to create comments in transaction",
					zap.Error(err),
				)
				return err
			}
		}

		global.Logger.Info("Posts and comments created successfully",
			zap.Uint("post_id", post.ID),
			zap.Int("comments_count", len(comments)),
		)
		return nil
	})
}
