package middleware

import (
	"crypto/subtle"
	"strings"

	"github.com/gin-gonic/gin"
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
			response.Unauthorized(c, "instance token is required")
			c.Abort()
			return
		}

		if subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) != 1 {
			response.Forbidden(c, "invalid instance token")
			c.Abort()
			return
		}

		c.Next()
	}
}
