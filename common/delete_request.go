package common

// DeleteRequest 删除请求包装
type DeleteRequest struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}
