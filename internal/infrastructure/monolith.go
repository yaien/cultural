package infrastructure

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gorilla/sessions"
	"github.com/robfig/cron/v3"
	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/lib/mail"
	"github.com/yaien/cultural/internal/lib/worker"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Monolith struct {
	Config          *Config
	GormDB          *gorm.DB
	MongoDB         *mongo.Database
	MongoClient     *mongo.Client
	Router          *http.ServeMux
	WebRouter       *http.ServeMux
	DashboardRouter *http.ServeMux
	SessionStore    sessions.Store
	Mail            mail.Mail
	StorageDriver   storage.Driver
	Queue           *worker.Queue
	Worker          *worker.Worker
	Cron            *cron.Cron
}

func NewMonolith() *Monolith {
	config := LoadConfig()

	var m Monolith
	m.Config = config
	m.GormDB = setupGormDB(config)
	m.SessionStore = setupSessionStore(config)
	m.Mail = setupMail(config)
	m.StorageDriver = setupStorage(config)
	m.Router = http.NewServeMux()
	m.WebRouter = http.NewServeMux()
	m.DashboardRouter = http.NewServeMux()
	m.Cron = cron.New()

	stream := worker.NewMemoryStream()
	store := worker.NewGormStore(m.GormDB)

	m.Queue = worker.NewQueue(store, stream)
	m.Worker = worker.New(store, stream)

	return &m

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

func setupStorage(config *Config) storage.Driver {
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

func setupGormDB(config *Config) *gorm.DB {

	option := &gorm.Config{
		Logger: logger.NewSlogLogger(slog.Default(), logger.Config{
			LogLevel:      logger.Info,
			SlowThreshold: 5 * time.Millisecond,
		}),
	}

	dialector := sqlite.Open(config.Sqlite.DSN)

	db, err := gorm.Open(dialector, option)
	if err != nil {
		log.Fatal("Failed to connect to sqlite database: ", err)
	}

	if err = db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		log.Fatal("Failed to enable foreign keys: ", err)
	}

	if err = db.Exec("PRAGMA journal_mode = WAL").Error; err != nil {
		log.Fatal("Failed to set journal mode: ", err)
	}

	return db
}
