package i18n

import "github.com/gin-gonic/gin"

var messages = map[string]map[string]string{
	"zh-CN": {
		"error.invalid_input":            "参数错误",
		"error.channel_not_found":        "通道不存在",
		"error.channel_id_required":      "通道 ID 不能为空",
		"error.invalid_channel_id":       "无效的通道 ID",
		"error.channel_verify_failed":    "通道验证失败",
		"error.device_not_found":        "设备不存在",
		"error.clipboard_not_found":      "剪贴板项不存在",
		"error.unauthorized":             "未授权",
		"error.forbidden":               "禁止访问",
		"error.instance_token_required": "实例令牌不能为空",
		"error.invalid_instance_token":  "无效的实例令牌",
		"error.internal_error":           "内部错误",
	},
	"en-US": {
		"error.invalid_input":            "Invalid input",
		"error.channel_not_found":        "Channel not found",
		"error.channel_id_required":     "Channel ID is required",
		"error.invalid_channel_id":       "Invalid channel ID",
		"error.channel_verify_failed":    "Channel verification failed",
		"error.device_not_found":         "Device not found",
		"error.clipboard_not_found":      "Clipboard item not found",
		"error.unauthorized":            "Unauthorized",
		"error.forbidden":                "Forbidden",
		"error.instance_token_required": "Instance token is required",
		"error.invalid_instance_token":  "Invalid instance token",
		"error.internal_error":           "Internal server error",
	},
}

// GetMessage 根据 message_key 和 Accept-Language 返回本地化消息
func GetMessage(ctx *gin.Context, key string) string {
	lang := getLang(ctx)
	if msg, ok := messages[lang][key]; ok {
		return msg
	}
	if msg, ok := messages["zh-CN"][key]; ok {
		return msg
	}
	return key
}

func getLang(ctx *gin.Context) string {
	lang := ctx.GetHeader("Accept-Language")
	if lang == "" {
		return "zh-CN"
	}
	if _, ok := messages[lang]; ok {
		return lang
	}
	return "zh-CN"
}