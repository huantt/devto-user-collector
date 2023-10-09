package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"user-collector/cmd"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("%s", r))
			slog.Info("Restart...")
			main()
		}
	}()

	var loggingLevel = new(slog.LevelVar)
	loggingLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: loggingLevel}))
	slog.SetDefault(logger)

	root := &cobra.Command{}
	root.AddCommand(cmd.CrawlUser("crawl-devto-users"))
	err := root.Execute()
	if err != nil {
		slog.Error(fmt.Sprintf("PANIC: %s", err.Error()))
		slog.Info("Restarting...")
		main()
	}
}
