package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// MySQLConfig MySQL数据库配置（可选）
type MySQLConfig struct {
	Host     string `yaml:"host,omitempty"`     // 数据库主机
	Port     int    `yaml:"port,omitempty"`     // 数据库端口
	Username string `yaml:"username,omitempty"` // 数据库用户名
	Password string `yaml:"password,omitempty"` // 数据库密码
	Database string `yaml:"database,omitempty"` // 数据库名称
	Charset  string `yaml:"charset,omitempty"`  // 字符集
}

// CORSConfig 控制浏览器跨域访问策略。
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins,omitempty"`
}

// SecurityConfig 存储通用安全限制。
type SecurityConfig struct {
	MaxBodyBytes int64 `yaml:"max_body_bytes,omitempty"`
}

// Config 存储应用程序配置
type Config struct {
	// 主机名，例如 "localhost" 或 "0.0.0.0"
	Host string `yaml:"host,omitempty"`
	// 端口号，例如 8080
	Port int `yaml:"port,omitempty"`
	// CORS配置（默认只允许本地前端开发地址跨域）
	CORS CORSConfig `yaml:"cors,omitempty"`
	// 安全配置
	Security SecurityConfig `yaml:"security,omitempty"`
	// MySQL配置（可选）
	MySQL *MySQLConfig `yaml:"mysql,omitempty"`
}

// 定义命令行参数
var (
	cmdPort = flag.Int("port", 0, "指定服务器端口，默认为8080")
	cmdP    = flag.Int("p", 0, "指定服务器端口的短参数形式")
)

// Load 加载应用程序配置
func Load() (*Config, error) {
	// 确保解析命令行参数
	if !flag.Parsed() {
		flag.Parse()
	}

	// 默认配置
	cfg := &Config{
		Host: "0.0.0.0",
		Port: 8080,
		CORS: CORSConfig{
			AllowedOrigins: []string{
				"http://localhost:3000",
				"http://127.0.0.1:3000",
			},
		},
		Security: SecurityConfig{
			MaxBodyBytes: 2 << 20, // 2 MiB
		},
	}

	// 应用命令行参数覆盖默认配置
	if *cmdPort > 0 {
		cfg.Port = *cmdPort
	} else if *cmdP > 0 {
		cfg.Port = *cmdP
	}

	applyEnvOverrides(cfg)
	normalize(cfg)

	return cfg, nil
}

// GetDatabaseType 根据配置确定使用的数据库类型
func (c *Config) GetDatabaseType() string {
	// 检查MySQL配置是否完整且有效
	if c.MySQL != nil && c.isValidMySQLConfig() {
		return "mysql"
	}
	return "sqlite"
}

// isValidMySQLConfig 检查MySQL配置是否完整有效
func (c *Config) isValidMySQLConfig() bool {
	if c.MySQL == nil {
		return false
	}

	// 检查必需的配置项
	if c.MySQL.Host == "" || c.MySQL.Username == "" || c.MySQL.Database == "" {
		return false
	}

	// 设置默认端口
	if c.MySQL.Port == 0 {
		c.MySQL.Port = 3306
	}

	// 设置默认字符集
	if c.MySQL.Charset == "" {
		c.MySQL.Charset = "utf8mb4"
	}

	return true
}

// GetDSN 根据数据库类型生成DSN连接字符串
func (c *Config) GetDSN() string {
	switch c.GetDatabaseType() {
	case "mysql":
		// MySQL DSN 格式: username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.MySQL.Username,
			c.MySQL.Password,
			c.MySQL.Host,
			c.MySQL.Port,
			c.MySQL.Database,
			c.MySQL.Charset,
		)
	case "sqlite":
		fallthrough
	default:
		// SQLite 使用默认路径
		homeDir, _ := os.UserHomeDir()
		appDir := filepath.Join(homeDir, ".cliplink")
		// #nosec G104 -- connection setup will fail later if the directory cannot be created.
		_ = os.MkdirAll(appDir, 0700)
		return filepath.Join(appDir, "cliplink.db")
	}
}

// GetServerAddress 获取服务器地址
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LoadFromFile 从指定文件路径加载配置
func LoadFromFile(configPath string) (*Config, error) {
	configPath = filepath.Clean(configPath)

	// 先获取默认配置
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	// 如果配置文件不存在，直接返回默认配置（SQLite）
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	// 读取配置文件
	// #nosec G304 -- configPath is intentionally supplied by CLI/container config.
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置文件（YAML格式）
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 命令行参数优先级高于配置文件
	if *cmdPort > 0 {
		cfg.Port = *cmdPort
	} else if *cmdP > 0 {
		cfg.Port = *cmdP
	}

	applyEnvOverrides(cfg)
	normalize(cfg)

	return cfg, nil
}

// SaveToFile 将配置保存到文件
func SaveToFile(cfg *Config, configPath string) error {
	configPath = filepath.Clean(configPath)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 转换为YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

func applyEnvOverrides(cfg *Config) {
	if origins := os.Getenv("CLIPLINK_ALLOWED_ORIGINS"); origins != "" {
		cfg.CORS.AllowedOrigins = splitCSV(origins)
	}

	if maxBodyBytes := os.Getenv("CLIPLINK_MAX_BODY_BYTES"); maxBodyBytes != "" {
		if value, err := strconv.ParseInt(maxBodyBytes, 10, 64); err == nil {
			cfg.Security.MaxBodyBytes = value
		}
	}
}

func normalize(cfg *Config) {
	cfg.CORS.AllowedOrigins = splitCSV(strings.Join(cfg.CORS.AllowedOrigins, ","))
	if cfg.Security.MaxBodyBytes <= 0 {
		cfg.Security.MaxBodyBytes = 2 << 20
	}
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
