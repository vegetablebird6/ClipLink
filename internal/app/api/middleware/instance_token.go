package middleware

import (
	"crypto/subtle"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/common/i18n"
	"github.com/xiaojiu/cliplink/internal/common/response"
)

const InstanceTokenHeader = "X-Instance-Token"

// InstanceTokenAuth protects instance-wide operations such as channel creation.
func InstanceTokenAuth(instanceToken string) gin.HandlerFunc {
	expected := strings.TrimSpace(instanceToken)

	return func(c *gin.Context) {
		if expected == "" {
			c.Next()
			return
		}

		actual := strings.TrimSpace(c.GetHeader(InstanceTokenHeader))
		if actual == "" {
			response.FailWithCode(c, 401, i18n.GetMessage(c, "error.instance_token_required"), "UNAUTHORIZED", "error.instance_token_required")
			c.Abort()
			return
		}

		if subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) != 1 {
			response.FailWithCode(c, 403, i18n.GetMessage(c, "error.forbidden"), "FORBIDDEN", "error.forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireInstanceTokenAuth protects destructive instance-wide operations.
func RequireInstanceTokenAuth(instanceToken string) gin.HandlerFunc {
	expected := strings.TrimSpace(instanceToken)

	return func(c *gin.Context) {
		if expected == "" {
			response.FailWithCode(c, 403, i18n.GetMessage(c, "error.forbidden"), "FORBIDDEN", "error.forbidden")
			c.Abort()
			return
		}

		actual := strings.TrimSpace(c.GetHeader(InstanceTokenHeader))
		if actual == "" {
			response.FailWithCode(c, 401, i18n.GetMessage(c, "error.instance_token_required"), "UNAUTHORIZED", "error.instance_token_required")
			c.Abort()
			return
		}

		if subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) != 1 {
			response.FailWithCode(c, 403, i18n.GetMessage(c, "error.forbidden"), "FORBIDDEN", "error.forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
