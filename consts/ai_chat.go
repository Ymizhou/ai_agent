package consts

type ChatModelType string

const (
	ChatModelTypeDeepSeek ChatModelType = "deepseek"
)

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleSystem    ChatRole = "system"
)

type ChatRespType string

const (
	ChatRespTypeGenerate ChatRespType = "generate"
	ChatRespTypeStream   ChatRespType = "stream"
)

type CodeGenarateType string

const (
	CodeGenarateTypeSingle CodeGenarateType = "single"
	CodeGenarateTypeMulti  CodeGenarateType = "multi"
)
