package util

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Host              string `mapstructure:"HOST"`
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	DBSourceTesting   string `mapstructure:"DB_SOURCE_TEST"`
	RedisAddress      string `mapstructure:"REDIS_ADDRESS"`
	RabbitMQAddress   string `mapstructure:"RABBITMQ_ADDRESS"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	GRPCSocketAddress string `mapstructure:"GRPC_SOCKET_ADDRESS"`

	SymmetricKey  string `mapstructure:"SYMMETRIC_KEY"`
	AdminUsername string `mapstructure:"ADMIN_USERNAME"`

	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	ApiPrefix          string `mapstructure:"API_PREFIX"`
	RegexIgnoreLogging string `mapstructure:"REGEX_IGNORE_LOGGING"`

	SocketUrl     string        `mapstructure:"SOCKET_ENDPOINT"`
	SocketTimeout time.Duration `mapstructure:"SOCKET__TIMEOUT"`

	CookieSameSite string `mapstructure:"COOKIE_SAME_SITE"`
	CookieSecure   bool   `mapstructure:"COOKIE_SECURE"`
	CookieUseHost  bool   `mapstructure:"COOKIE_USE_HOST"`

	DefaultAuthenticationUsername string `mapstructure:"DEFAULT_AUTHENTICATION_USERNAME"`

	AccessControlAllowOrigin string `mapstructure:"ACCESS_CONTROL_ALLOW_ORIGIN"`

	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

var Configs = Config{}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		LoadConfig(path)
	})
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&Configs)
	return &Configs, err
}
