package infrastructure

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/yaien/cultural/internal/library/mail"
	"github.com/yaien/cultural/internal/library/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Monolith struct {
	Config       *Config
	MongoDB      *mongo.Database
	MongoClient  *mongo.Client
	Router       *http.ServeMux
	WebRouter    *http.ServeMux
	SessionStore sessions.Store
	Mail         mail.Mail
	Storage      storage.Storage
}

func NewMonolith() *Monolith {
	config := LoadConfig()

	setupOauthProviders(config)

	return &Monolith{
		Config:       config,
		MongoDB:      setupMongoDB(config),
		SessionStore: setupSessionStore(config),
		Mail:         setupMail(config),
		Storage:      setupStorage(config),
		Router:       http.NewServeMux(),
		WebRouter:    http.NewServeMux(),
	}
}

func setupOauthProviders(config *Config) {
	store := sessions.NewCookieStore([]byte(config.SessionConfig.Secret))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode
	store.Options.Secure = store.Options.Secure

	gothic.Store = store

	goth.UseProviders(
		google.New(config.Google.ClientID, config.Google.ClientSecret, config.Server.URL+"/auth/google/callback", "email", "profile"),
	)
}

func setupMail(config *Config) mail.Mail {
	switch config.Mail.Provider {
	case "mailtrap":
		client, err := mail.NewMailtrap(mail.MailtrapOptions{
			Token:   config.Mail.Mailtrap.Token,
			Sandbox: config.Mail.Mailtrap.Sandbox,
			InboxID: config.Mail.Mailtrap.InboxID,
		})
		if err != nil {
			log.Fatal("Failed to create mailtrap client: ", err)
		}
		return client
	default:
		log.Fatal("Unsupported mail provider:", config.Mail.Provider)
		return nil
	}
}

func setupStorage(config *Config) storage.Storage {
	switch config.Storage.Provider {
	case "local":
		return storage.NewLocal(config.Storage.Local.Path)
	default:
		log.Fatal("Unsupported storage provider:", config.Storage.Provider)
		return nil
	}
}

func setupSessionStore(config *Config) sessions.Store {
	store := sessions.NewCookieStore([]byte(config.SessionConfig.Secret))
	store.Options.Secure = config.SessionConfig.Secure
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode
	store.Options.MaxAge = int((7 * 24 * time.Hour).Seconds()) // 7 days
	return store
}

func setupMongoDB(config *Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(config.MongoDB.URI)
	opts.SetLoggerOptions(options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug).SetSink(&sink{}))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("Failed to connect to mongodb:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping mongodb:", err)
	}

	log.Println("Connected to mongodb successfully")

	database := client.Database(config.MongoDB.Database)

	return database
}

var _ options.LogSink = (*sink)(nil)

type sink struct {
}

func (s *sink) Info(level int, msg string, args ...any) {
	slog.With(args...).Debug("Mongodb", "level", level, "msg", msg)
}

func (s *sink) Error(err error, msg string, args ...any) {
	slog.With(args...).Error("Mongodb", "error", err, "msg", msg)
}
