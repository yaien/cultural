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
	SessionConfig SessionConfig
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

func NewConfig() *Config {
	return &Config{
		Debug: viper.GetBool("DEBUG"),
		Server: ServerConfig{
			Addr: viper.GetString("SERVER_ADDR"),
			URL:  viper.GetString("SERVER_URL"),
		},
		SessionConfig: SessionConfig{
			Key:    viper.GetString("SESSION_KEY"),
			Secure: viper.GetBool("SESSION_SECURE"),
		},
		MongoDB: MongoDBConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Database: viper.GetString("MONGODB_DATABASE"),
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
