package main

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/yaien/cultural/internal/infrastructure"
)

func main() {
	root := cmd()
	root.AddCommand(serve())

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
			app := infrastructure.NewMonolith()
			err := register(app)
			if err != nil {
				log.Fatal("Failed to register modules:", err)
			}

			log.Printf("MongoDB Database: %s", app.Config.MongoDB.Database)
			log.Printf("MongoDB URI: %s", app.Config.MongoDB.URI)
			log.Printf("App is running on %s", app.Config.Server.URL)

			err = http.ListenAndServe(app.Config.Server.Addr, nil)
			if err != nil {
				log.Fatal("Failed to start server:", err)
			}
		},
	}

	return cmd
}
