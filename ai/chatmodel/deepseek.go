package chatmodel

import (
	"context"
	"errors"

	"aicode/config"
	"aicode/consts"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/model"
)

func initDeepSeek(cfg *config.Config) error {
	// 校验参数
	deepSeekCfg := cfg.AI.DeepSeek
	if deepSeekCfg.APIKey == "" {
		return errors.New("深度求索 API 密钥不能为空")
	}
	if deepSeekCfg.Model == "" {
		return errors.New("深度求索 模型不能为空")
	}
	if deepSeekCfg.BaseURL == "" {
		return errors.New("深度求索 BaseURL 不能为空")
	}

	// 注册聊天模型工厂函数
	registerChatModel(string(consts.ChatModelTypeDeepSeek), func(ctx context.Context) (model.BaseChatModel, error) {
		return deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			APIKey:  cfg.AI.DeepSeek.APIKey,
			Model:   cfg.AI.DeepSeek.Model,
			BaseURL: cfg.AI.DeepSeek.BaseURL,
		})
	})
	return nil
}
