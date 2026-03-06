package chatmodel

import (
	"context"
	"fmt"

	"aicode/config"
	"aicode/consts"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type ChatModelFactory func(ctx context.Context) (model.BaseChatModel, error)

var chatModelRegistry = make(map[string]ChatModelFactory)

func InitChatModel(cfg *config.Config) (map[string]ChatModelFactory, error) {
	err := initDeepSeek(cfg)
	if err != nil {
		return nil, err
	}
	return chatModelRegistry, nil
}

// registerChatModel 注册聊天模型进入工厂
func registerChatModel(name string, factory ChatModelFactory) {
	chatModelRegistry[name] = factory
}

func GetChatModel(ctx context.Context, name string) (model.BaseChatModel, error) {
	create, ok := chatModelRegistry[name]
	if !ok {
		return nil, fmt.Errorf("不支持的模型类型: %s", name)
	}

	return create(ctx)
}

func AutoChat(ctx context.Context,
	model model.BaseChatModel, messages []*schema.Message, respType consts.ChatRespType) (
	*schema.Message, *schema.StreamReader[*schema.Message], error) {
	resp, err := GetGenarateRespType(ctx, model, respType)
	if err != nil {
		return nil, nil, err
	}
	return resp.DoChat(ctx, messages)
}
