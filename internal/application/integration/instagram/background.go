package instagram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yaien/cultural/internal/lib/primitive"
	"github.com/yaien/cultural/internal/lib/worker"
	"golang.org/x/oauth2"
)

const TaskName = "sync-instagram"

type TaskData struct {
	OrganizationID primitive.ID `json:"organization_id"`
}

func (i *Instagram) RegisterBackgroundProcess(c *cron.Cron, q *worker.Queue, w *worker.Worker) {

	_, _ = c.AddFunc("12 21 * * *", i.cron(q))

	w.Register(worker.H{
		Name:       TaskName,
		MaxRetries: 3,
		Handler:    worker.HandlerFunc[TaskData](i.handle),
	})
}

func (i *Instagram) cron(q *worker.Queue) cron.FuncJob {
	return func() {
		slog.Info("starting instagram sync background process")

		ctx := context.Background()

		integrations, err := i.integrations.Where("name = ?", i.Name()).Find(ctx)
		if err != nil {
			slog.Error("failed getting integrations by name", "error", err)
			return
		}

		for _, integration := range integrations {
			if !integration.Data.Connected {
				continue

			}

			task := worker.Task{
				Name: TaskName,
				Data: TaskData{
					OrganizationID: integration.OrganizationID,
				},
			}

			if err := q.Push(ctx, task); err != nil {
				slog.Error("failed pushing task to queue", "error", err)
			}
		}

		slog.Info("finished pushing instagram sync tasks to queue")
	}
}

func (i *Instagram) handle(ctx context.Context, data *TaskData) error {
	integration, err := i.integrations.Where("organization_id = ? and name = ?", data.OrganizationID, i.Name()).First(ctx)
	if err != nil {
		return fmt.Errorf("failed getting integration by organization ID and name: %w", err)
	}

	config, err := i.configs.GetByOrganizationID(ctx, data.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed getting config by organization ID: %w", err)
	}

	auth := i.OAuthConfig(&config)

	token := &oauth2.Token{AccessToken: integration.Data.Token}

	client := NewClient(token, auth)

	if integration.Data.ExpireAt.Add(-48 * time.Hour).Before(time.Now()) {
		token, err = client.RefreshLongToken(ctx)
		if err != nil {
			return fmt.Errorf("failed getting long token")
		}

		client.SetToken(token)
	}

	user, err := client.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("failed getting user: %w", err)
	}

	posts, err := client.GetPosts(ctx)
	if err != nil {
		return fmt.Errorf("failed getting posts: %w", err)
	}

	if len(posts) > 10 {
		posts = posts[:6]
	}

	instagramData := Data{
		Connected: true,
		User:      user,
		Posts:     posts,
		Token:     token.AccessToken,
		ExpireAt:  time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}

	if err := i.Save(ctx, config.OrganizationID, instagramData); err != nil {
		return fmt.Errorf("failed saving data: %w", err)
	}

	return nil

}
