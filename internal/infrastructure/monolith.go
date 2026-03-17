package infrastructure

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/robfig/cron/v3"
	"github.com/yaien/cultural/internal/library/mail"
	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/library/worker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

type OAuth struct {
	Google    *oauth2.Config
	Instagram *oauth2.Config
}

type Monolith struct {
	Config          *Config
	MongoDB         *mongo.Database
	MongoClient     *mongo.Client
	Router          *http.ServeMux
	WebRouter       *http.ServeMux
	DashboardRouter *http.ServeMux
	SessionStore    sessions.Store
	Mail            mail.Mail
	Storage         storage.Storage
	Queue           *worker.Queue
	Worker          *worker.Worker
	Cron            *cron.Cron
	OAuth           *OAuth
}

func NewMonolith() *Monolith {
	config := LoadConfig()

	setupOauthProviders(config)

	var m Monolith
	m.Config = config
	m.MongoDB = setupMongoDB(config)
	m.SessionStore = setupSessionStore(config)
	m.Mail = setupMail(config)
	m.Storage = setupStorage(config)
	m.Router = http.NewServeMux()
	m.WebRouter = http.NewServeMux()
	m.DashboardRouter = http.NewServeMux()
	m.Cron = cron.New()

	stream := worker.NewMemoryStream()
	store := worker.NewMongoStore(m.MongoDB, "")

	m.Queue = worker.NewQueue(store, stream)
	m.Worker = worker.New(store, stream)

	return &m

}

func setupOauthProviders(config *Config) {
	store := sessions.NewCookieStore([]byte(config.Session.Secret))
	store.Options.Path = ""
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode
	store.Options.Secure = config.Session.Secure

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
	store := sessions.NewCookieStore([]byte(config.Session.Secret))
	store.Options.Secure = config.Session.Secure
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
		slog.Error("Failed to connect to mongodb:", "error", err)
		os.Exit(1)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error("Failed to ping mongodb:", "error", err)
		os.Exit(1)
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
