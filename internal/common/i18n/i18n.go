package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var messages = map[string]map[string]string{
	"zh-CN": {
		"error.invalid_input":           "参数错误",
		"error.channel_not_found":       "通道不存在",
		"error.channel_id_required":     "通道 ID 不能为空",
		"error.invalid_channel_id":      "无效的通道 ID",
		"error.channel_verify_failed":   "通道验证失败",
		"error.device_not_found":        "设备不存在",
		"error.clipboard_not_found":     "剪贴板项不存在",
		"error.unauthorized":            "未授权",
		"error.forbidden":               "禁止访问",
		"error.not_found":               "资源不存在",
		"error.instance_token_required": "实例令牌不能为空",
		"error.invalid_instance_token":  "无效的实例令牌",
		"error.internal_error":          "内部错误",
	},
	"en-US": {
		"error.invalid_input":           "Invalid input",
		"error.channel_not_found":       "Channel not found",
		"error.channel_id_required":     "Channel ID is required",
		"error.invalid_channel_id":      "Invalid channel ID",
		"error.channel_verify_failed":   "Channel verification failed",
		"error.device_not_found":        "Device not found",
		"error.clipboard_not_found":     "Clipboard item not found",
		"error.unauthorized":            "Unauthorized",
		"error.forbidden":               "Forbidden",
		"error.not_found":               "Resource not found",
		"error.instance_token_required": "Instance token is required",
		"error.invalid_instance_token":  "Invalid instance token",
		"error.internal_error":          "Internal server error",
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
	h := ctx.GetHeader("Accept-Language")
	if h == "" {
		return "zh-CN"
	}
	// 支持带权重格式如 "en-US,en;q=0.9"，取第一个
	if idx := strings.Index(h, ","); idx >= 0 {
		h = strings.TrimSpace(h[:idx])
	}
	h = strings.TrimSpace(h)
	if _, ok := messages[h]; ok {
		return h
	}
	// 前缀匹配：支持 "en" 匹配 "en-US"
	for lang := range messages {
		if strings.HasPrefix(lang, h) {
			return lang
		}
	}
	return "zh-CN"
}
