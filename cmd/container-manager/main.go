package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"container-manager/internal/logger"

	"github.com/urfave/cli/v3"
)

var (
	appName    = "sentinel"
	version    = "local"
	commitHash = "unknown"
	buildDate  = "unknown"
)

func getBuildVersion() string {
	return fmt.Sprintf(
		"\nrelease: %s\ncommit hash: %s\nbuild date: %s\ngo version: %s",
		version,
		commitHash,
		buildDate,
		runtime.Version(),
	)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	defer cancel()

	l := logger.New(logger.Config{
		Level:  slog.LevelInfo,
		Format: "json", // можно поменять на "text" для более читаемого формата
		ExtraAttrs: map[string]any{
			"service": appName,
			"version": version,
			"build":   buildDate,
			"commit":  commitHash,
		},
	})

	app := &cli.Command{
		Name:    appName,
		Usage:   "A CLI application for " + appName,
		Version: version,
		Suggest: true,
		Commands: []*cli.Command{
			startCMD(l),
			versionCMD(),
		},
	}

	// run cli runner
	if err := app.Run(ctx, os.Args); err != nil {
		l.Fatalf("failed to run cli runner: %s", err)
	}
}
