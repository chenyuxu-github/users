package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config 是项目的完整配置，对应 config/config.toml 中的顶层配置段。
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig 是 HTTP 服务配置。
type ServerConfig struct {
	// Host 是服务监听的主机地址，例如 0.0.0.0 或 127.0.0.1。
	Host string `mapstructure:"host"`
	// Port 是服务监听端口。
	Port int `mapstructure:"port"`
}

// DatabaseConfig 是 MySQL 数据库配置。
type DatabaseConfig struct {
	// Host 是数据库主机地址。
	Host string `mapstructure:"host"`
	// Port 是数据库端口。
	Port int `mapstructure:"port"`
	// Username 是数据库用户名。
	Username string `mapstructure:"username"`
	// Password 是数据库密码。
	Password string `mapstructure:"password"`
	// Name 是要连接的数据库名。
	Name string `mapstructure:"name"`
	// Charset 是连接使用的字符集。
	Charset string `mapstructure:"charset"`
	// ParseTime 控制 MySQL 时间类型是否自动解析为 Go 的 time.Time。
	ParseTime bool `mapstructure:"parse_time"`
	// Loc 是数据库时间的时区配置，例如 Local。
	Loc string `mapstructure:"loc"`
}

// Load 使用 Viper 读取配置文件，并允许通过环境变量覆盖配置项。
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	// 将 server.port 这类配置键映射为 SERVER_PORT 环境变量名。
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnv(v,
		"server.host",
		"server.port",
		"database.host",
		"database.port",
		"database.username",
		"database.password",
		"database.name",
		"database.charset",
		"database.parse_time",
		"database.loc",
	)

	var cfg Config
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	// 对必填配置做基础校验，启动时尽早暴露配置问题。
	if cfg.Server.Port <= 0 {
		return nil, fmt.Errorf("server.port is required")
	}
	if cfg.Database.Host == "" {
		return nil, fmt.Errorf("database.host is required")
	}
	if cfg.Database.Port <= 0 {
		return nil, fmt.Errorf("database.port is required")
	}
	if cfg.Database.Username == "" {
		return nil, fmt.Errorf("database.username is required")
	}
	if cfg.Database.Name == "" {
		return nil, fmt.Errorf("database.name is required")
	}

	return &cfg, nil
}

// Address 返回 Gin 启动服务需要的监听地址，例如 0.0.0.0:8080。
func (c ServerConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// DSN 根据拆分后的数据库配置组装 MySQL 连接字符串。
func (c DatabaseConfig) DSN() string {
	params := url.Values{}
	if c.Charset != "" {
		params.Set("charset", c.Charset)
	}
	params.Set("parseTime", strconv.FormatBool(c.ParseTime))
	if c.Loc != "" {
		params.Set("loc", c.Loc)
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?%s",
		c.Username,
		c.Password,
		net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		c.Name,
		params.Encode(),
	)
}

// bindEnv 显式绑定允许被环境变量覆盖的配置键。
func bindEnv(v *viper.Viper, keys ...string) {
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}
