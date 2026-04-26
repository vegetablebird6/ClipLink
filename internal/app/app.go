package app

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/app/api/routes"
	"github.com/xiaojiu/cliplink/internal/app/usecase"
	"github.com/xiaojiu/cliplink/internal/config"
	"github.com/xiaojiu/cliplink/internal/infra/db"
	"github.com/xiaojiu/cliplink/internal/infra/persistence"
)

// BuildRouter 初始化所有依赖并返回 gin.Engine
func BuildRouter() (*gin.Engine, error) {
	return BuildRouterWithConfig(nil)
}

// BuildRouterWithConfig 使用指定配置初始化所有依赖并返回 gin.Engine
func BuildRouterWithConfig(cfg *config.Config) (*gin.Engine, error) {
	// 1. 加载配置
	if cfg == nil {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return nil, fmt.Errorf("加载配置失败: %w", err)
		}
	}

	// 2. 初始化数据库
	if _, err := db.InitWithConfig(cfg); err != nil {
		// 如果配置的数据库初始化失败且当前是MySQL，无感切换到SQLite
		if cfg.GetDatabaseType() == "mysql" {
			// 创建SQLite配置
			sqliteConfig := &config.Config{
				Host: cfg.Host,
				Port: cfg.Port,
				// MySQL设为nil，自动使用SQLite
				MySQL: nil,
			}

			if _, fallbackErr := db.InitWithConfig(sqliteConfig); fallbackErr != nil {
				return nil, fmt.Errorf("数据库初始化失败: %w", fallbackErr)
			}
		} else {
			return nil, fmt.Errorf("数据库初始化失败: %w", err)
		}
	}

	// 3. 创建 gin 引擎
	router := gin.Default()

	// 4. 设置安全中间件（必须在注册路由前 use）
	router.Use(SecurityHeaders())
	router.Use(RequestBodyLimit(cfg.Security.MaxBodyBytes))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.CORS.AllowedOrigins
	corsConfig.AllowCredentials = false
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "X-Channel-ID"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	if len(corsConfig.AllowOrigins) > 0 {
		router.Use(cors.New(corsConfig))
	}

	// 5. 创建仓库
	channelRepo := persistence.NewChannelRepository()
	clipboardRepo := persistence.NewClipboardRepository()
	deviceRepo := persistence.NewDeviceRepository()
	syncHistoryRepo := persistence.NewSyncHistoryRepository()

	// 6. 创建服务
	channelService := usecase.NewChannelService(channelRepo, clipboardRepo, deviceRepo)
	clipboardService := usecase.NewClipboardService(clipboardRepo, syncHistoryRepo)
	deviceService := usecase.NewDeviceService(deviceRepo)
	statsService := usecase.NewStatsService(deviceRepo, clipboardRepo, channelRepo, syncHistoryRepo)
	syncService := usecase.NewSyncService(syncHistoryRepo)

	// 7. 注册 API 路由
	routes.SetupRouter(
		router,
		channelService,
		clipboardService,
		deviceService,
		statsService,
		syncService,
	)

	return router, nil
}

// SecurityHeaders adds conservative browser security headers for the UI and API.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Writer.Header()
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("X-Frame-Options", "DENY")
		headers.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		headers.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), clipboard-read=(self), clipboard-write=(self)")

		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			headers.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

// RequestBodyLimit prevents oversized payloads from exhausting memory or storage.
func RequestBodyLimit(maxBodyBytes int64) gin.HandlerFunc {
	if maxBodyBytes <= 0 {
		maxBodyBytes = 2 << 20
	}

	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)
		c.Next()
	}
}
