package response

import "awesomeProject/homework04/common/types"

// PostResponse 文章响应
type PostResponse struct {
	ID        uint            `json:"id"`
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	UserId    uint            `json:"userId"`
	CreatedAt types.LocalTime `json:"createdAt"`
	UpdatedAt types.LocalTime `json:"updatedAt"`
}

// PostListResponse 文章列表响应（简化版）
type PostListResponse struct {
	ID        uint            `json:"id"`
	Title     string          `json:"title"`
	UserId    uint            `json:"userId"`
	CreatedAt types.LocalTime `json:"createdAt"`
}
