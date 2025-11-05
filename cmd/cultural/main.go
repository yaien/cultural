package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"github.com/yaien/cultural/internal/modules/configs"
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/library/render"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"github.com/yaien/cultural/internal/modules/landing/interface/web/assets"
)

func main() {
	root := cmd()
	root.AddCommand(serve())
	root.AddCommand(migrate())
	root.AddCommand(revert())
	root.AddCommand(invite())
	root.AddCommand(edit())

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

func cmd() *cobra.Command {
	return &cobra.Command{
		Use: "cultural",
	}
}

func serve() *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			mono := infrastructure.NewMonolith()
			err := register(mono)
			if err != nil {
				log.Fatal("Failed to register modules:", err)
			}

			log.Printf("MongoDB Database: %s", mono.Config.MongoDB.Database)
			log.Printf("MongoDB URI: %s", mono.Config.MongoDB.URI)
			log.Printf("App is running on %s", mono.Config.Server.URL)

			err = http.ListenAndServe(mono.Config.Server.Addr, mono.Router)
			if err != nil {
				log.Fatal("Failed to start server:", err)
			}
		},
	}

	return cmd
}

func migrate() *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			mono := infrastructure.NewMonolith()
			err := migrations.Migrate(cmd.Context(), mono.MongoDB)
			if err != nil {
				log.Fatal("Failed to run migrations:", err)
			}
			log.Println("Migrations applied successfully")
		},
	}

	return cmd
}

func revert() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "revert [migration name]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			mono := infrastructure.NewMonolith()
			err := migrations.Revert(cmd.Context(), args[0], mono.MongoDB)
			if err != nil {
				log.Fatal("Failed to revert migrations:", err)
			}
			log.Println("Migrations reverted successfully")
		},
	}

	return cmd
}

func invite() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "invite [email] [display-name]",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			mono := infrastructure.NewMonolith()
			err := register(mono)
			if err != nil {
				log.Fatal("Failed to register modules:", err)
			}

			cfg, err := infrastructure.Resolve(mono, &configs.Module{})
			if err != nil {
				log.Fatal("Failed to resolve configs module:", err)
			}

			ctx := cmd.Context()

			config, err := cfg.App.GetConfigByHost(ctx, mono.Config.Init.Host)
			if err != nil {
				log.Fatal("Failed to get config by host:", err)
			}

			_, err = cfg.App.CreateInvitation(ctx, &commands.CreateInvitationRequest{
				OrganizationID:  config.OrganizationID,
				UserEmail:       args[0],
				UserDisplayName: args[1],
				RolePermissions: []string{"*"},
				RoleName:        "Admin",
				ExpiresAt:       time.Now().Add(3 * time.Hour),
			})

			if err != nil {
				log.Fatal("Failed to create invitation:", err)
			}

			log.Printf("Invitation sent successfully")
		},
	}

	return cmd

}

func edit() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit [host] [name]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			host, name := args[0], args[1]

			mono := infrastructure.NewMonolith()
			err := register(mono)
			if err != nil {
				return fmt.Errorf("failed to register modules: %w", err)
			}

			cfg, err := infrastructure.Resolve(mono, &configs.Module{})
			if err != nil {
				return fmt.Errorf("failed to resolve configs module: %w", err)
			}

			config, err := cfg.App.GetConfigByHost(cmd.Context(), host)
			if err != nil {
				return fmt.Errorf("failed to get config by host: %w", err)
			}

			page, ok := config.Pages[name]
			if !ok {
				return fmt.Errorf("page with name %s not found", name)
			}

			file, err := os.CreateTemp(".", "template.*.json")
			if err != nil {
				return fmt.Errorf("failed to create temp file: %w", err)
			}

			stat, err := file.Stat()
			if err != nil {
				return fmt.Errorf("failed to stat temp file: %w", err)
			}

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "    ")
			err = encoder.Encode(page.Body)
			if err != nil {
				return fmt.Errorf("failed to write page body to temp file: %w", err)
			}

			err = file.Close()
			if err != nil {
				return fmt.Errorf("failed to close temp file: %w", err)
			}

			defer func() {
				err := os.Remove(file.Name())
				if err != nil {
					slog.Error("failed to remove temp file:", "err", err)
				}
			}()

			changes := make(chan struct{}, 1000)
			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return fmt.Errorf("failed to create file watcher: %w", err)
			}

			defer watcher.Close()

			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							continue
						}

						if event.Name == stat.Name() {
							changes <- struct{}{}
						}

					case err, ok := <-watcher.Errors:
						if ok {
							slog.Error(err.Error())
						}

					}
				}
			}()

			err = watcher.Add("./")
			if err != nil {
				return fmt.Errorf("failed to add watcher to file: %w", err)
			}

			router := http.NewServeMux()

			router.Handle("GET /assets/static/landing/", http.StripPrefix("/assets/static/landing/", http.FileServerFS(assets.FS)))

			router.HandleFunc("GET /assets/landing/styles.css", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/css")
				err := models.Styles.Execute(w, config)
				if err != nil {
					http.Error(w, "Failed to generate styles", http.StatusInternalServerError)
					return
				}
			})

			router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
				data, err := os.ReadFile(file.Name())
				if err != nil {
					http.Error(w, "Failed to read template file", http.StatusInternalServerError)
					return
				}

				var body models.Node
				err = json.Unmarshal(data, &body)
				if err != nil {
					http.Error(w, "Failed to unmarshal template JSON", http.StatusInternalServerError)
					return
				}

				page.Body = body

				ctx := context.WithValue(r.Context(), models.ConfigContextKey, config)

				w.Header().Set("Content-Type", "text/html")
				_ = render.Page(page, nil, render.WithSSE("/sse")).Render(ctx, w)
			})

			router.HandleFunc("GET /sse", func(w http.ResponseWriter, r *http.Request) {
				flusher, ok := w.(http.Flusher)
				if !ok {
					http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")
				w.WriteHeader(http.StatusOK)

				flusher.Flush()

				ctx := context.WithValue(r.Context(), models.ConfigContextKey, config)

				for {
					select {
					case <-changes:
						data, err := os.ReadFile(file.Name())
						if err != nil {
							return
						}

						var body models.Node
						err = json.Unmarshal(data, &body)
						if err != nil {
							return
						}

						var buff bytes.Buffer
						_ = render.Render(body, nil).Render(ctx, &buff)

						_, _ = fmt.Fprintf(w, "event: changes\ndata: %s\n\n", strings.ReplaceAll(buff.String(), "\n", "\\n"))
						flusher.Flush()
					case <-r.Context().Done():
						return
					}
				}
			})

			server := httptest.NewServer(router)

			err = browser.OpenURL(server.URL)
			if err != nil {
				slog.Error("failed to open browser:", "err", err)
			}

			<-interrupt

			return nil
		},
	}

	return cmd
}
