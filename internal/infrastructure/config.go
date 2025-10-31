package infrastructure

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
)

type Config struct {
	Debug         bool
	MongoDB       MongoDBConfig
	Server        ServerConfig
	Init          InitConfig
	SessionConfig SessionConfig
	SMTPConfig    SMTPConfig
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type ServerConfig struct {
	Addr string
	URL  string
}

type SessionConfig struct {
	Key    string
	Secure bool
}

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type InitConfig struct {
	Host  string
	Url   string
	Title string
}

func LoadConfig() *Config {
	return &Config{
		Debug: viper.GetBool("DEBUG"),
		Server: ServerConfig{
			Addr: viper.GetString("SERVER_ADDR"),
			URL:  viper.GetString("SERVER_URL"),
		},
		Init: InitConfig{
			Host:  viper.GetString("INIT_HOST"),
			Url:   viper.GetString("INIT_URL"),
			Title: viper.GetString("INIT_TITLE"),
		},
		SessionConfig: SessionConfig{
			Key:    viper.GetString("SESSION_KEY"),
			Secure: viper.GetBool("SESSION_SECURE"),
		},
		MongoDB: MongoDBConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Database: viper.GetString("MONGODB_DATABASE"),
		},
		SMTPConfig: SMTPConfig{
			Host: viper.GetString("SMTP_HOST"),
			Port: viper.GetInt("SMTP_PORT"),
			User: viper.GetString("SMTP_USER"),
			Pass: viper.GetString("SMTP_PASS"),
		},
	}
}

func init() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	if viper.GetBool("DEBUG") {
		slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
		})))

	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})))
	}

}
