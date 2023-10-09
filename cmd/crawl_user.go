package cmd

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"user-collector/handler/devto"
	"user-collector/pkg/forem"
	"user-collector/repository/repository"
)

func CrawlUser(use string) *cobra.Command {
	var from, to, concurrent int
	var proxy string

	command := &cobra.Command{
		Use: use,
		Run: func(cmd *cobra.Command, args []string) {
			const tableName = "devto_users"
			if proxy != "" {
				slog.Info("Using proxy: " + proxy)
			}
			db, err := sql.Open("sqlite3", "data/db/database.sqlite")
			if err != nil {
				slog.Error(err.Error())
				os.Exit(1)
			}
			defer db.Close()
			userRepo := repository.NewUserRepository(db, tableName)
			devtoService := forem.NewService(forem.DevToEndpoint, concurrent, proxy)
			userCollector := devto.NewUserCollector(userRepo, devtoService)
			userCollector.Collect(context.Background(), concurrent, from, to)
		},
	}
	command.Flags().IntVar(&from, "from", 1, "User ID from")
	command.Flags().IntVar(&to, "to", 1, "User ID to")
	command.Flags().IntVar(&concurrent, "concurrent", 2, "Concurrent requests")
	command.Flags().StringVar(&proxy, "proxy", "", "Proxy")

	if err := command.MarkFlagRequired("from"); err != nil {
		panic(err)
	}
	if err := command.MarkFlagRequired("to"); err != nil {
		panic(err)
	}
	return command
}
