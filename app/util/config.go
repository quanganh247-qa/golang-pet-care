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
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`

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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	GoogleAPIKey       string `mapstructure:"GOOGLE_API_KEY"`

	OpenFDAAPIKey string `mapstructure:"OPENFDA_API_KEY"`
<<<<<<< HEAD
=======
	GoogleAPIKey       string `mapstructure:"GOOGLE_API_KEY"`
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
	GoogleAPIKey       string `mapstructure:"GOOGLE_API_KEY"`
>>>>>>> e859654 (Elastic search)

	GoongAPIKey  string `mapstructure:"GOONG_API_KEY"`
	GoongBaseURL string `mapstructure:"GOONG_BASE_URL"`

	VietQRBaseURL   string `mapstructure:"VIETQR_BASE_URL"`
	VietQRAPIKey    string `mapstructure:"VIETQR_API_KEY"`
	VietQRClientKey string `mapstructure:"VIETQR_CLIENT_KEY"`
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD

	PaypalClientID     string `mapstructure:"PAYPAL_CLIENT_ID"`
	PaypalClientSecret string `mapstructure:"PAYPAL_CLIENT_SECRET"`
	PaypalEnvironment  string `mapstructure:"PAYPAL_ENVIRONMENT"`
	PaypalURL          string `mapstructure:"PAYPAL_URL"`

	ElasticsearchURL string `mapstructure:"ELASTICSEARCH_URL"`
	MinIOEndpoint    string `mapstructure:"MINIO_ENDPOINT"`
	MinIOAccessKey   string `mapstructure:"MINIO_ACCESS_KEY"`
	MinIOSecretKey   string `mapstructure:"MINIO_SECRET_KEY"`
	MinIOSSL         bool   `mapstructure:"MINIO_SSL"`
<<<<<<< HEAD
<<<<<<< HEAD
=======

	GoongAPIKey  string `mapstructure:"GOONG_API_KEY"`
	GoongBaseURL string `mapstructure:"GOONG_BASE_URL"`
>>>>>>> 4625843 (added goong maps api)
=======
>>>>>>> c449ffc (feat: cart api)
=======

	PaypalClientID     string `mapstructure:"PAYPAL_CLIENT_ID"`
	PaypalClientSecret string `mapstructure:"PAYPAL_CLIENT_SECRET"`
	PaypalEnvironment  string `mapstructure:"PAYPAL_ENVIRONMENT"`
	PaypalURL          string `mapstructure:"PAYPAL_URL"`

	ElasticsearchURL string `mapstructure:"ELASTICSEARCH_URL"`
	MinIOEndpoint    string `mapstructure:"MINIO_ENDPOINT"`
	MinIOAccessKey   string `mapstructure:"MINIO_ACCESS_KEY"`
	MinIOSecretKey   string `mapstructure:"MINIO_SECRET_KEY"`
	MinIOSSL         bool   `mapstructure:"MINIO_SSL"`
>>>>>>> e859654 (Elastic search)
=======

	ClickatellAPIKey string `mapstructure:"CLICKATELL_API_KEY"`
	ClickatellAPIID  string `mapstructure:"CLICKATELL_API_ID"`
	ClickatellURL    string `mapstructure:"CLICKATELL_URL"`
>>>>>>> 4ccd381 (Update appointment flow)
=======
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)
=======

	GoongAPIKey  string `mapstructure:"GOONG_API_KEY"`
	GoongBaseURL string `mapstructure:"GOONG_BASE_URL"`
>>>>>>> 4625843 (added goong maps api)
=======
>>>>>>> c449ffc (feat: cart api)
=======

	PaypalClientID     string `mapstructure:"PAYPAL_CLIENT_ID"`
	PaypalClientSecret string `mapstructure:"PAYPAL_CLIENT_SECRET"`
	PaypalEnvironment  string `mapstructure:"PAYPAL_ENVIRONMENT"`
	PaypalURL          string `mapstructure:"PAYPAL_URL"`

	ElasticsearchURL string `mapstructure:"ELASTICSEARCH_URL"`
	MinIOEndpoint    string `mapstructure:"MINIO_ENDPOINT"`
	MinIOAccessKey   string `mapstructure:"MINIO_ACCESS_KEY"`
	MinIOSecretKey   string `mapstructure:"MINIO_SECRET_KEY"`
	MinIOSSL         bool   `mapstructure:"MINIO_SSL"`
>>>>>>> e859654 (Elastic search)
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
