package request

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=100"`
	Content string `json:"content" binding:"required,min=1,max=10000"`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=100"`
	Content string `json:"content" binding:"omitempty,min=1,max=10000"`
}
