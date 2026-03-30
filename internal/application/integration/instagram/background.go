package instagram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yaien/cultural/internal/lib/worker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

const TaskName = "sync-instagram"

func (i *Instagram) RegisterBackgroundProcess(c *cron.Cron, q *worker.Queue, w *worker.Worker) {

	c.AddFunc("12 21 * * *", func() {
		slog.Info("starting instagram sync background process")

		ctx := context.Background()

		integrations, err := i.integrations.GetByName(ctx, i.Name())
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
				Data: map[string]any{"organizationId": integration.OrganizationID},
			}

			if err := q.Push(ctx, task); err != nil {
				slog.Error("failed pushing task to queue", "error", err)
			}
		}

		slog.Info("finished pushing instagram sync tasks to queue")
	})

	w.Register(worker.H{
		Name:       TaskName,
		MaxRetries: 3,
		Handler: worker.HandlerFunc(func(ctx context.Context, taskData map[string]any) error {
			organizationID, ok := taskData["organizationId"].(primitive.ObjectID)
			if !ok {
				slog.Error("failed asserting organizationId to ObjectID")
				return nil
			}

			integration, err := i.integrations.GetByOrganizationIDAndName(ctx, organizationID, i.Name())
			if err != nil {
				return fmt.Errorf("failed getting integration by organization ID and name: %w", err)
			}

			config, err := i.configs.GetByOrganizationID(ctx, organizationID)
			if err != nil {
				return fmt.Errorf("failed getting config by organization ID: %w", err)
			}

			auth := i.OAuthConfig(config)

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

		}),
	})
}
