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
	Google        GoogleConfig
	Storage       StorageConfig
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
	Secret string
	Secure bool
}

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	APIKey       string
}

type InitConfig struct {
	Host  string
	Url   string
	Title string
	Email string
}

type StorageConfig struct {
	Provider string
	Local    Local
}

type Local struct {
	Path string
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
			Email: viper.GetString("INIT_EMAIL"),
		},
		SessionConfig: SessionConfig{
			Secret: viper.GetString("SESSION_SECRET"),
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
		Google: GoogleConfig{
			ClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
			ClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
			APIKey:       viper.GetString("GOOGLE_API_KEY"),
		},
		Storage: StorageConfig{
			Provider: viper.GetString("STORAGE_PROVIDER"),
			Local: Local{
				Path: viper.GetString("STORAGE_LOCAL_PATH"),
			},
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
