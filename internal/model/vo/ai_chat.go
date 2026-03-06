package vo

import "aicode/consts"

// AIChatRequest 聊天测试请求结构
type AIChatRequest struct {
	Model    consts.ChatModelType `json:"model" binding:"required"`
	Question string               `json:"question" binding:"required"`
	History  []AIChatMessage      `json:"history,omitempty"`
}

// AIChatMessage 聊天消息结构（用于前端传递）
type AIChatMessage struct {
	Role    consts.ChatRole `json:"role"` // "user" 或 "assistant"
	Content string          `json:"content"`
}

// AIChatResponse 聊天测试响应结构
type AIChatResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
