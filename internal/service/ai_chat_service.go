package service

import (
	"aicode/internal/model/vo"
	"context"

	"github.com/cloudwego/eino/schema"
)

type AIChatService interface {
	ChatGenerate(ctx context.Context, params *vo.AIChatRequest) (string, error)
	ChatStream(ctx context.Context, params *vo.AIChatRequest) (*schema.StreamReader[*schema.Message], error)
}
