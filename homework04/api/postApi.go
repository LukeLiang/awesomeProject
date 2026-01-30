package api

import (
	"errors"
	"strconv"

	"awesomeProject/homework04/common"
	"awesomeProject/homework04/common/request"
	"awesomeProject/homework04/common/response"
	"awesomeProject/homework04/global"
	"awesomeProject/homework04/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostApi struct {
}

func (a *PostApi) Create(c *gin.Context) {
	var req request.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Warn("Create post: invalid params",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		response.FailWithMessage(c, response.ErrInvalidParams, err.Error())
		return
	}

	currentUser := c.MustGet(common.CurrentUserPrefix).(common.GlobalUser)

	post := &model.Posts{
		Title:   req.Title,
		Content: req.Content,
		UserId:  currentUser.ID,
	}

	if err := postService.Create(post); err != nil {
		global.Logger.Error("Create post: service error",
			zap.Error(err),
			zap.Uint("user_id", currentUser.ID),
		)
		response.Fail(c, response.ErrPostCreate)
		return
	}

	global.Logger.Info("Posts created",
		zap.Uint("post_id", post.ID),
		zap.Uint("user_id", currentUser.ID),
	)
	response.SuccessWithData(c, "文章创建成功", response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserId:    post.UserId,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

func (a *PostApi) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.Logger.Warn("Update post: invalid id",
			zap.String("id", c.Param("id")),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	// 获取现有文章
	post, err := postService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, response.ErrPostNotFound)
		} else {
			response.Fail(c, response.ErrDBQuery)
		}
		return
	}

	// 权限校验
	currentUser := c.MustGet(common.CurrentUserPrefix).(common.GlobalUser)
	if post.UserId != currentUser.ID {
		global.Logger.Warn("Update post: permission denied",
			zap.Uint("post_id", post.ID),
			zap.Uint("owner_id", post.UserId),
			zap.Uint("request_user_id", currentUser.ID),
		)
		response.Fail(c, response.ErrPermissionDenied)
		return
	}

	// 绑定更新参数
	var req request.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Warn("Update post: invalid params",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		response.FailWithMessage(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := postService.Update(post); err != nil {
		global.Logger.Error("Update post: service error",
			zap.Error(err),
			zap.Uint("post_id", post.ID),
		)
		response.Fail(c, response.ErrPostUpdate)
		return
	}

	global.Logger.Info("Posts updated",
		zap.Uint("post_id", post.ID),
		zap.Uint("user_id", currentUser.ID),
	)
	response.SuccessWithMessage(c, "文章更新成功")
}

func (a *PostApi) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.Logger.Warn("Get post detail: invalid id",
			zap.String("id", c.Param("id")),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	post, err := postService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, response.ErrPostNotFound)
		} else {
			response.Fail(c, response.ErrDBQuery)
		}
		return
	}

	response.Success(c, response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserId:    post.UserId,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

func (a *PostApi) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		global.Logger.Warn("Delete post: invalid id",
			zap.String("id", c.Param("id")),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	// 获取现有文章
	post, err := postService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, response.ErrPostNotFound)
		} else {
			response.Fail(c, response.ErrDBQuery)
		}
		return
	}

	// 权限校验
	currentUser := c.MustGet(common.CurrentUserPrefix).(common.GlobalUser)
	if post.UserId != currentUser.ID {
		global.Logger.Warn("Delete post: permission denied",
			zap.Uint("post_id", post.ID),
			zap.Uint("owner_id", post.UserId),
			zap.Uint("request_user_id", currentUser.ID),
		)
		response.Fail(c, response.ErrPermissionDenied)
		return
	}

	if err := postService.Delete(uint(id)); err != nil {
		global.Logger.Error("Delete post: service error",
			zap.Error(err),
			zap.Uint("post_id", uint(id)),
		)
		response.Fail(c, response.ErrPostDelete)
		return
	}

	global.Logger.Info("Posts deleted",
		zap.Uint("post_id", uint(id)),
		zap.Uint("user_id", currentUser.ID),
	)
	response.SuccessWithMessage(c, "文章删除成功")
}

func (a *PostApi) ListPosts(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		global.Logger.Warn("ListPosts: invalid userId",
			zap.String("userId", userIdStr),
			zap.String("client_ip", c.ClientIP()),
		)
		response.Fail(c, response.ErrInvalidParams)
		return
	}

	// 绑定分页参数
	var pageReq request.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		global.Logger.Warn("ListPosts: invalid page params",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		response.FailWithMessage(c, response.ErrInvalidParams, err.Error())
		return
	}

	posts, total, err := postService.List(uint(userId), &pageReq)
	if err != nil {
		global.Logger.Error("ListPosts: query failed",
			zap.Uint64("userId", userId),
			zap.Error(err),
		)
		response.Fail(c, response.ErrDBQuery)
		return
	}

	// 转换为响应 VO
	postList := make([]response.PostListResponse, len(posts))
	for i, p := range posts {
		postList[i] = response.PostListResponse{
			ID:        p.ID,
			Title:     p.Title,
			UserId:    p.UserId,
			CreatedAt: p.CreatedAt,
		}
	}

	global.Logger.Debug("ListPosts: success",
		zap.Uint64("userId", userId),
		zap.Int64("total", total),
	)
	response.Success(c, response.NewPageResponse(postList, total, pageReq.GetPage(), pageReq.GetPageSize()))
}
