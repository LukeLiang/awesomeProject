package request

// PageRequest 分页请求
type PageRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100"`
}

// GetPage 获取页码，默认为1
func (p *PageRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页大小，默认为10
func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	return p.PageSize
}

// GetOffset 获取偏移量
func (p *PageRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}
