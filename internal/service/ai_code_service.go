package service

import (
	"context"

	"aicode/internal/model/vo"
)

type AICodeService interface {
	// CodeGenerateStream 启动流式代码生成，结果逐块写入 ch，stream 读取与文件写入均在内部 goroutine 中异步完成
	CodeGenerateStream(ctx context.Context, params *vo.AICodeRequest, ch chan<- vo.CodeStreamResult) error
	CodeGenerate(ctx context.Context, params *vo.AICodeRequest) (string, error)
}
