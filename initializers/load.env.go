package initializers

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ClientOrigin   string        `mapstructure:"CLIENT_ORIGIN"`
	EmailFrom      string        `mapstructure:"EMAIL_FROM"`
	DBUserPassword string        `mapstructure:"MYSQL_PASSWORD"`
	DBName         string        `mapstructure:"MYSQL_DB"`
	DBPort         string        `mapstructure:"MYSQL_PORT"`
	ServerPort     string        `mapstructure:"PORT"`
	DBUserName     string        `mapstructure:"MYSQL_USER"`
	SMTPUser       string        `mapstructure:"SMTP_USER"`
	DBHost         string        `mapstructure:"MYSQL_HOST"`
	SMTPPass       string        `mapstructure:"SMTP_PASS"`
	TokenSecret    string        `mapstructure:"TOKEN_SECRET"`
	SMTPHost       string        `mapstructure:"SMTP_HOST"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`
	SMTPPort       int           `mapstructure:"SMTP_PORT"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRED_IN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
