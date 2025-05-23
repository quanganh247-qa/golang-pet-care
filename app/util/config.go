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
	DATABASE_URL      string `mapstructure:"DATABASE_URL"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	RedisAddress      string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
	RedisUsername     string `mapstructure:"REDIS_USERNAME"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	DebugMode         bool   `mapstructure:"DEBUG_MODE"`

	SymmetricKey  string `mapstructure:"SYMMETRIC_KEY"`
	AdminUsername string `mapstructure:"ADMIN_USERNAME"`

	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	ApiPrefix string `mapstructure:"API_PREFIX"`

	CookieSameSite string `mapstructure:"COOKIE_SAME_SITE"`
	CookieSecure   bool   `mapstructure:"COOKIE_SECURE"`
	CookieUseHost  bool   `mapstructure:"COOKIE_USE_HOST"`

	AccessControlAllowOrigin string `mapstructure:"ACCESS_CONTROL_ALLOW_ORIGIN"`

	EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`

	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
	GoogleAPIKey       string `mapstructure:"GOOGLE_API_KEY"`

	OpenFDAAPIKey string `mapstructure:"OPENFDA_API_KEY"`
	OpenAIAPIKey  string `mapstructure:"OPENAI_API_KEY"`

	GoongAPIKey  string `mapstructure:"GOONG_API_KEY"`
	GoongBaseURL string `mapstructure:"GOONG_BASE_URL"`

	VietQRBaseURL   string `mapstructure:"VIETQR_BASE_URL"`
	VietQRAPIKey    string `mapstructure:"VIETQR_API_KEY"`
	VietQRClientKey string `mapstructure:"VIETQR_CLIENT_KEY"`

	RoboflowAPIKey string `mapstructure:"ROBOFLOW_API_KEY"`

	MinIOEndpoint  string `mapstructure:"MINIO_ENPOINT"`
	MinIOAccessKey string `mapstructure:"MINIO_USERNAME"`
	MinIOSecretKey string `mapstructure:"MINIO_PASSWORD"`
	MinIOSSL       bool   `mapstructure:"MINIO_SSL"`
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
