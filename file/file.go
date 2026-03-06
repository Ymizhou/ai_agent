package file

import (
	"aicode/config"
	"aicode/consts"
	"aicode/internal/model/vo"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// buildAppDir 构造应用目录路径：基础路径 + "app" + appId
func buildAppDir(appId string) string {
	cfg := config.GetConfig()
	return filepath.Join(cfg.File.StoreBasePath, "app/"+appId)
}

// writeFile 确保目录存在后将内容写入指定文件
func writeFile(filePath, content string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败 [%s]: %w", filePath, err)
	}
	return nil
}

// StoreToSingalFile 将单文件模式生成结果存储到本地
// content 必须是符合 SingleHTMLResult 结构的 JSON 字符串
// 写入文件：{basePath}/app{appId}/index.html
func StoreToSingalFile(ctx context.Context, appId, content string) error {
	var result vo.SingleHTMLResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return fmt.Errorf("解析单文件生成结果失败: %w", err)
	}
	dir := buildAppDir(appId)
	return writeFile(filepath.Join(dir, "index.html"), result.HTML)
}

// StoreToMultiFile 将多文件模式生成结果存储到本地
// content 必须是符合 MultiHTMLResult 结构的 JSON 字符串
// 写入文件：{basePath}/app{appId}/index.html、style.css、script.js
func StoreToMultiFile(ctx context.Context, appId, content string) error {
	var result vo.MultiHTMLResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return fmt.Errorf("解析多文件生成结果失败: %w", err)
	}
	dir := buildAppDir(appId)
	fileMap := map[string]string{
		"index.html": result.HTML,
		"style.css":  result.CSS,
		"script.js":  result.JavaScript,
	}
	for name, fileContent := range fileMap {
		if err := writeFile(filepath.Join(dir, name), fileContent); err != nil {
			return err
		}
	}
	return nil
}

// stripMarkdownCodeBlock 剥离 AI 模型可能附加的 markdown 代码块包装（```json ... ```）
// 若内容不含代码块标记则原样返回
var markdownCodeBlockRe = regexp.MustCompile("(?s)^```[a-zA-Z]*\\n(.+?)\\n?```\\s*$")

func stripMarkdownCodeBlock(content string) string {
	content = strings.TrimSpace(content)
	if m := markdownCodeBlockRe.FindStringSubmatch(content); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return content
}

// StoreByGenType 根据代码生成类型自动分派到对应存储函数，供外部统一调用
func StoreByGenType(ctx context.Context,
	genType consts.CodeGenarateType, appId, content string) error {
	if appId == "" {
		return fmt.Errorf("appId 不能为空")
	}
	if content == "" {
		return fmt.Errorf("content 不能为空")
	}
	// 剥离 AI 返回内容中可能存在的 markdown 代码块包装
	content = stripMarkdownCodeBlock(content)
	switch genType {
	case consts.CodeGenarateTypeSingle:
		return StoreToSingalFile(ctx, appId, content)
	case consts.CodeGenarateTypeMulti:
		return StoreToMultiFile(ctx, appId, content)
	default:
		return fmt.Errorf("不支持的代码生成类型: %s", genType)
	}
}
