package impl

import (
	"aicode/ai/chatmodel"
	"aicode/config"
	"aicode/consts"
	"aicode/file"
	"aicode/internal/model/vo"
	"aicode/internal/service"
	"context"
	"io"
	"os"
	"strings"

	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

// AICodeServiceImpl ai代码服务实现
type AICodeServiceImpl struct {
}

// NewAICodeService 创建ai代码服务实例
func NewAICodeService() service.AICodeService {
	return &AICodeServiceImpl{}
}

func (s *AICodeServiceImpl) CodeGenerateStream(ctx context.Context,
	params *vo.AICodeRequest, ch chan<- vo.CodeStreamResult) error {
	// 构建消息列表
	messages := dealCodeMessages(ctx, params.GenType, params.Question, params.History)

	// 获取模型实例
	chat, _ := chatmodel.GetChatModel(ctx, string(params.Model))
	// 调用模型，获取 stream
	_, stream, err := chatmodel.AutoChat(ctx,
		chat, messages, consts.ChatRespTypeStream)
	if err != nil {
		return err
	}

	// 异步读取 stream，逐块写入 channel；全量内容收集完毕后写入文件
	go func() {
		defer close(ch)
		defer recover()

		var buf strings.Builder
		for {
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					// stream 正常结束，将全量内容写入文件
					if storeErr := file.StoreByGenType(ctx,
						params.GenType, params.AppId, buf.String()); storeErr != nil {
						logrus.Errorf("写入文件失败: %v", storeErr)
						ch <- vo.CodeStreamResult{Err: storeErr}
					}
					return
				}
				// stream 读取出错
				ch <- vo.CodeStreamResult{Err: err}
				return
			}

			if msg != nil && msg.Content != "" {
				buf.WriteString(msg.Content)
				ch <- vo.CodeStreamResult{Content: msg.Content}
			}
		}
	}()

	return nil
}

func (s *AICodeServiceImpl) CodeGenerate(ctx context.Context,
	params *vo.AICodeRequest) (string, error) {
	// 构建消息列表
	messages := dealCodeMessages(ctx, params.GenType, params.Question, params.History)

	// 获取模型实例
	chat, _ := chatmodel.GetChatModel(ctx, string(params.Model))
	// 调用模型
	message, _, err := chatmodel.AutoChat(ctx,
		chat, messages, consts.ChatRespTypeGenerate)
	if err != nil {
		return "", err
	}
	// 文件存储
	err = file.StoreByGenType(ctx,
		params.GenType, params.AppId, message.Content)
	if err != nil {
		return "", err
	}
	return message.Content, nil
}

func dealCodeMessages(ctx context.Context,
	genType consts.CodeGenarateType, question string,
	history []vo.AICodeMessage) []*schema.Message {
	messages := make([]*schema.Message, 0)
	// 添加历史对话
	if len(history) > 0 {
		for _, msg := range history {
			switch msg.Role {
			case consts.ChatRoleSystem:
				messages = append(messages,
					schema.SystemMessage(msg.Content))
			case consts.ChatRoleUser:
				messages = append(messages,
					schema.UserMessage(msg.Content))
			case consts.ChatRoleAssistant:
				messages = append(messages,
					schema.AssistantMessage(msg.Content, []schema.ToolCall{}))
			}
		}
	} else {
		// 添加默认系统提示词
		systemPrompt := getSystemPrompt(genType)
		messages = append(messages, schema.SystemMessage(systemPrompt))
	}

	// 添加当前问题
	messages = append(messages, schema.UserMessage(question))
	return messages
}

func getSystemPrompt(genType consts.CodeGenarateType) string {
	dirPath := ""
	cfg := config.GetConfig()
	switch genType {
	case consts.CodeGenarateTypeSingle:
		dirPath = cfg.AI.SystemPromptDir.SingalGenerate
	case consts.CodeGenarateTypeMulti:
		dirPath = cfg.AI.SystemPromptDir.MultiGenerate
	default:
		return ""
	}
	bytes, err := os.ReadFile(dirPath)
	if err != nil {
		logrus.Errorf("读取系统提示词文件失败: %s", err.Error())
		return ""
	}
	content := string(bytes)
	return content
}
