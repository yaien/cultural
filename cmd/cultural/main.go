package main

import (
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/infrastructure/migrations"
	"github.com/yaien/cultural/internal/modules/configs"
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
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
		Use:  "invite [email]",
		Args: cobra.ExactArgs(1),
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
			email := args[0]

			config, err := cfg.App.GetConfigByHost(ctx, mono.Config.Init.Host)
			if err != nil {
				log.Fatal("Failed to get config by host:", err)
			}

			_, err = cfg.App.CreateInvitation(ctx, &commands.CreateInvitationRequest{
				OrganizationID: config.OrganizationID,
				BaseURL:        config.Url,
				Email:          email,
				Permissions:    []string{"*"},
				ExpiresAt:      time.Now().Add(3 * time.Hour),
				Name:           "Admin",
			})

			if err != nil {
				log.Fatal("Failed to create invitation:", err)
			}

			log.Printf("Invitation sent to %s successfully", email)
		},
	}

	return cmd

}
