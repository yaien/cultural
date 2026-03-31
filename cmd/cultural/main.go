package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"github.com/yaien/cultural/internal/application"
	"github.com/yaien/cultural/internal/application/admin"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"github.com/yaien/cultural/internal/web"
)

func main() {
	root := cmd()
	root.AddCommand(serve())
	root.AddCommand(migrate())
	root.AddCommand(revert())
	root.AddCommand(invite())

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

			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
			defer stop()

			mono := infrastructure.NewMonolith()
			app := application.New(mono)
			web.Register(mono, app)

			log.Printf("App is running on %s", mono.Config.Server.URL)

			go func() {

				if err := migrations.Migrate(ctx, mono.GormDB); err != nil {
					slog.Error("Failed running migrations", "error", err)
					return
				}

				log.Println("Migrations checked successfully")

				mono.Cron.Start()
				log.Println("Cron started successfully")

				mono.Worker.Start()
				log.Println("Worker started successfully")

				<-ctx.Done()
				if err := ctx.Err(); err != nil {
					slog.Info("context done", "error", context.Cause(ctx))
				}

				mono.Worker.Stop()
				log.Println("Worker stopped successfully")

				<-mono.Cron.Stop().Done()
				log.Println("Cron stopped successfully")

				os.Exit(0)
			}()

			if mono.Config.Server.TLS {
				err := http.ListenAndServeTLS(
					mono.Config.Server.Addr,
					mono.Config.Server.CertFile,
					mono.Config.Server.KeyFile,
					mono.Router,
				)

				if err != nil {
					log.Fatal("Failed to start server:", err)
				}

				return
			}

			err := http.ListenAndServe(mono.Config.Server.Addr, mono.Router)
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
			err := migrations.Migrate(cmd.Context(), mono.GormDB)
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
			err := migrations.Revert(cmd.Context(), args[0], mono.GormDB)
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
			app := application.New(mono)
			ctx := cmd.Context()

			config, err := app.Label.Configs.GetByHost(ctx, mono.Config.Init.Host)
			if err != nil {
				log.Fatal("Failed to get config by host:", err)
			}

			opts := &admin.CreateInvitationOptions{
				OrganizationID:  config.OrganizationID,
				UserEmail:       args[0],
				UserDisplayName: args[1],
				RolePermissions: []string{"*"},
				RoleName:        "Admin",
				ExpiresAt:       time.Now().Add(3 * time.Hour),
			}

			if _, err = app.Admin.Invitations.Create(ctx, opts); err != nil {
				log.Fatal("Failed to create invitation:", err)
			}

			log.Printf("Invitation sent successfully")
		},
	}

	return cmd

}
