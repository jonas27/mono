package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	exitCodeErr  = 1
	appName      = "hourly"
	fileLocation = "hourly/data.csv"
)

func fullpath(path string) (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return fmt.Sprintf("%s/%s", dirname, path), nil
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	args := os.Args

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	app := &cli.App{
		Name: appName,
		Commands: []*cli.Command{
			{
				Name:  "login",
				Usage: "Writes a login to the hourly log file.",
				Action: func(c *cli.Context) error {
					return run(ctx, log, c)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Usage:   "The file to write to.",
						EnvVars: []string{"HOURLY_FILE"},
						Value:   fileLocation,
					},
				},
			}, {
				Name:  "logout",
				Usage: "Writes a logout to the hourly log file.",
				Action: func(c *cli.Context) error {
					return run(ctx, log, c)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Usage:   "The file to write to.",
						EnvVars: []string{"HOURLY_FILE"},
						Value:   fileLocation,
					},
				},
			},
		},
	}

	if err := app.Run(args); err != nil {
		cancel()
		fmt.Fprintf(os.Stderr, "%s stopped with error: %v\n", appName, err)
		os.Exit(exitCodeErr)
	}
}

func run(ctx context.Context, log *slog.Logger, c *cli.Context) error {
	log.Debug("running command", "command", c.Command.Name)
	path := c.String("file")
	if path == fileLocation {
		var err error
		path, err = fullpath(fileLocation)
		if err != nil {
			return fmt.Errorf("failed to get full path: %w", err)
		}
	}

	if err := ensureDir(log, path); err != nil {
		return fmt.Errorf("failed to ensure directory exists: %w", err)
	}
	if err := ensureFile(log, path); err != nil {
		return fmt.Errorf("failed to ensure file exists: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	writeToFile(c.Command.Name, f)
	return nil
}

func ensureDir(log *slog.Logger, file string) error {
	log.Debug("ensure dir", "file", file)
	dir := filepath.Dir(file)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Debug("create dir", "file", file)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	log.Debug("dir and file exists", "dir", dir, "file", fileLocation)
	return nil
}

func ensureFile(log *slog.Logger, file string) error {
	log.Debug("ensure file", "file", file)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Debug("create file")
		if _, err = os.Create(file); err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	}
	log.Debug("file exists", "file", file)
	return nil
}

func writeToFile(method string, file *os.File) error {
	time := time.Now().Unix()
	line := fmt.Sprintf("%s,%d\n", method, time)
	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}
