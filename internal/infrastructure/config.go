package infrastructure

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
)

type Config struct {
	MongoDB MongoDBConfig
	Server  ServerConfig
	Init    InitConfig
	Session SessionConfig
	Google  GoogleConfig
	Storage StorageConfig
	Mail    MailConfig
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
	Local    LocalStorageConfig
}

type LocalStorageConfig struct {
	Path string
}

type MailConfig struct {
	Provider string
	Mailtrap MailtrapMailConfig
}

type MailtrapMailConfig struct {
	Token   string
	Sandbox bool
	InboxID string
}

func LoadConfig() *Config {
	return &Config{
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
		Session: SessionConfig{
			Secret: viper.GetString("SESSION_SECRET"),
			Secure: viper.GetBool("SESSION_SECURE"),
		},
		MongoDB: MongoDBConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Database: viper.GetString("MONGODB_DATABASE"),
		},
		Mail: MailConfig{
			Provider: viper.GetString("MAIL_PROVIDER"),
			Mailtrap: MailtrapMailConfig{
				Token:   viper.GetString("MAIL_MAILTRAP_TOKEN"),
				Sandbox: viper.GetBool("MAIL_MAILTRAP_SANDBOX"),
				InboxID: viper.GetString("MAIL_MAILTRAP_INBOX_ID"),
			},
		},
		Google: GoogleConfig{
			ClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
			ClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
			APIKey:       viper.GetString("GOOGLE_API_KEY"),
		},
		Storage: StorageConfig{
			Provider: viper.GetString("STORAGE_PROVIDER"),
			Local: LocalStorageConfig{
				Path: viper.GetString("STORAGE_LOCAL_PATH"),
			},
		},
	}
}

func init() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	logLevel := slog.Level(viper.GetInt("LOG_LEVEL"))
	switch viper.GetString("LOG_FORMAT") {
	case "tint":
		slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			TimeFormat: time.DateTime,
			Level:      logLevel,
		})))

	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: logLevel,
		})))

	default:
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: logLevel,
		})))
	}

}
