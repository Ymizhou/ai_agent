package common

// PageRequest 分页请求封装
type PageRequest struct {
	PageNum   int    `json:"pageNum" form:"pageNum"`     // 当前页号
	PageSize  int    `json:"pageSize" form:"pageSize"`   // 页面大小
	SortField string `json:"sortField" form:"sortField"` // 排序字段
	SortOrder string `json:"sortOrder" form:"sortOrder"` // 排序顺序（默认降序）
}

// DefaultPageRequest 返回默认分页参数
func DefaultPageRequest() PageRequest {
	return PageRequest{
		PageNum:   1,
		PageSize:  10,
		SortOrder: "descend",
	}
}
