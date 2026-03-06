package chatmodel

import (
	"aicode/consts"
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type GenarateRespType interface {
	DoChat(ctx context.Context, input []*schema.Message) (
		*schema.Message, *schema.StreamReader[*schema.Message], error)
}

func GetGenarateRespType(ctx context.Context,
	chatModel model.BaseChatModel,
	respType consts.ChatRespType) (GenarateRespType, error) {
	switch respType {
	case consts.ChatRespTypeGenerate:
		return &GenarateAll{ChatModel: chatModel}, nil
	case consts.ChatRespTypeStream:
		return &GenarateStream{ChatModel: chatModel}, nil
	default:
		return nil, fmt.Errorf("不支持的响应类型: %s", respType)
	}
}

type GenarateAll struct {
	ChatModel model.BaseChatModel
}

func (g *GenarateAll) DoChat(ctx context.Context,
	messages []*schema.Message) (*schema.Message, *schema.StreamReader[*schema.Message], error) {
	message, err := g.ChatModel.Generate(ctx, messages)
	if err != nil {
		return nil, nil, err
	}
	return message, nil, nil
}

type GenarateStream struct {
	ChatModel model.BaseChatModel
}

func (g *GenarateStream) DoChat(ctx context.Context,
	messages []*schema.Message) (*schema.Message, *schema.StreamReader[*schema.Message], error) {
	stream, err := g.ChatModel.Stream(ctx, messages)
	if err != nil {
		return nil, nil, err
	}
	return nil, stream, nil
}
