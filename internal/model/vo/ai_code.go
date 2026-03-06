package vo

import "aicode/consts"

// AIChatRequest 聊天测试请求结构
type AICodeRequest struct {
	AppId    string                  `json:"appId" binding:"required"`
	Model    consts.ChatModelType    `json:"model" binding:"required"`
	GenType  consts.CodeGenarateType `json:"genType" binding:"required"`
	Question string                  `json:"question" binding:"required"`
	History  []AICodeMessage         `json:"history,omitempty"`
}

// AICodeMessage 代码消息结构（用于前端传递）
type AICodeMessage struct {
	Role    consts.ChatRole `json:"role"` // "user" 或 "assistant"
	Content string          `json:"content"`
}

// CodeStreamResult channel 中传递的流式结果单元
type CodeStreamResult struct {
	Content string
	Err     error
}
