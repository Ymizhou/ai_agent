package vo

// SingleHTMLResult 单文件模式生成结果，对应 singal_html_generate 提示词输出的 JSON 结构
type SingleHTMLResult struct {
	HTML string `json:"html"`
}

// MultiHTMLResult 多文件模式生成结果，对应 multi_html_generate 提示词输出的 JSON 结构
type MultiHTMLResult struct {
	HTML       string `json:"html"`
	CSS        string `json:"css"`
	JavaScript string `json:"javascript"`
}
