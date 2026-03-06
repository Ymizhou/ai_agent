package impl

import (
	"aicode/ai/chatmodel"
	"aicode/consts"
	"aicode/internal/model/vo"
	"aicode/internal/service"
	"context"

	"github.com/cloudwego/eino/schema"
)

// AIChatServiceImpl ai聊天服务实现
type AIChatServiceImpl struct {
}

// NewAIChatService 创建ai聊天服务实例
func NewAIChatService() service.AIChatService {
	return &AIChatServiceImpl{}
}

func (s *AIChatServiceImpl) ChatGenerate(ctx context.Context,
	params *vo.AIChatRequest) (string, error) {
	// 构建消息列表
	messages := dealChatMessages(params.Question, params.History)

	// 获取模型实例
	chat, _ := chatmodel.GetChatModel(ctx, string(params.Model))
	// 调用模型
	message, _, err := chatmodel.AutoChat(ctx,
		chat, messages, consts.ChatRespTypeGenerate)
	if err != nil {
		return "", err
	}
	return message.Content, nil
}

func (s *AIChatServiceImpl) ChatStream(ctx context.Context,
	params *vo.AIChatRequest) (*schema.StreamReader[*schema.Message], error) {
	// 构建消息列表
	messages := dealChatMessages(params.Question, params.History)

	// 获取模型实例
	chat, _ := chatmodel.GetChatModel(ctx, string(params.Model))
	// 调用模型
	_, stream, err := chatmodel.AutoChat(ctx, chat, messages, consts.ChatRespTypeStream)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func dealChatMessages(question string, history []vo.AIChatMessage) []*schema.Message {
	messages := make([]*schema.Message, 0)
	// 添加历史对话
	for _, msg := range history {
		switch msg.Role {
		case consts.ChatRoleUser:
			messages = append(messages,
				schema.UserMessage(msg.Content))
		case consts.ChatRoleAssistant:
			messages = append(messages,
				schema.AssistantMessage(msg.Content, []schema.ToolCall{}))
		}
	}
	// 添加当前问题
	messages = append(messages, schema.UserMessage(question))
	return messages
}
