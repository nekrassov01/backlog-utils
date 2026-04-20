package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/nekrassov01/backlog-utils/backlog"
	"github.com/nekrassov01/backlog-utils/backlog/wiki"
	"github.com/nekrassov01/backlog-utils/log"
	"github.com/nekrassov01/backlog-utils/version"
	"github.com/urfave/cli/v3"
)

const name = "bkl"

var logger = &slog.Logger{}

func newCmd(w, ew io.Writer) *cli.Command {
	logger = log.NewLogger(ew, slog.LevelInfo.String())

	loglevel := &cli.StringFlag{
		Name:    "log-level",
		Usage:   "set log level",
		Sources: cli.EnvVars("BACKLOG_LOG_LEVEL"),
		Value:   slog.LevelInfo.String(),
	}

	baseURL := &cli.StringFlag{
		Name:    "base-url",
		Usage:   "set backlog base url",
		Sources: cli.EnvVars("BACKLOG_URL"),
	}

	apiKey := &cli.StringFlag{
		Name:    "api-key",
		Usage:   "set backlog api key",
		Sources: cli.EnvVars("BACKLOG_API_KEY"),
	}

	projectKey := &cli.StringFlag{
		Name:     "project-key",
		Usage:    "set backlog project key",
		Required: true,
	}

	pattern := &cli.StringFlag{
		Name:  "pattern",
		Usage: "set pattern to search for wiki pages",
	}

	wikiID := &cli.IntFlag{
		Name:     "wiki-id",
		Usage:    "set backlog wiki id",
		Required: true,
	}

	oldString := &cli.StringFlag{
		Name:     "old",
		Usage:    "set string to be replaced in wiki page",
		Required: true,
	}

	newString := &cli.StringFlag{
		Name:     "new",
		Usage:    "set new string after replacement in wiki page",
		Required: true,
	}

	pairs := &cli.StringSliceFlag{
		Name:     "pairs",
		Usage:    "set pairs of old and new repalacements for wiki page",
		Required: true,
	}

	beforeWiki := func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		logger = log.NewLogger(cmd.Writer, cmd.String(loglevel.Name))

		transport := backlog.NewRetryableTransport(1*time.Second, 30*time.Second, 5, 3000)
		client, err := wiki.NewClient(
			cmd.String(baseURL.Name),
			cmd.String(apiKey.Name),
			backlog.WithWriter(cmd.Writer),
			backlog.WithTransport(transport),
		)
		if err != nil {
			return nil, err
		}

		cmd.Metadata["client"] = client
		return ctx, nil
	}

	listWiki := func(_ context.Context, cmd *cli.Command) error {
		logger.Info("started")

		client := cmd.Metadata["client"].(*wiki.Client)
		pages, err := client.List(cmd.String(projectKey.Name), cmd.String(pattern.Name))
		if err != nil {
			return err
		}

		enc := json.NewEncoder(cmd.Writer)
		for _, page := range pages {
			if err := enc.Encode(page); err != nil {
				return err
			}
		}

		logger.Info("stopped")
		return nil
	}

	renameWiki := func(_ context.Context, cmd *cli.Command) error {
		logger.Info("started")

		client := cmd.Metadata["client"].(*wiki.Client)
		page, err := client.Get(cmd.Int(wikiID.Name))
		if err != nil {
			return err
		}

		if err := client.Rename(page, cmd.String(oldString.Name), cmd.String(newString.Name)); err != nil {
			return err
		}

		logger.Info("stopped")
		return nil
	}

	replaceWiki := func(_ context.Context, cmd *cli.Command) error {
		logger.Info("started")

		client := cmd.Metadata["client"].(*wiki.Client)
		pages, err := client.List(cmd.String(projectKey.Name), cmd.String(pattern.Name))
		if err != nil {
			return err
		}

		for _, page := range pages {
			if err := client.Replace(page, cmd.StringSlice(pairs.Name)...); err != nil {
				return err
			}
		}

		logger.Info("stopped")
		return nil
	}

	renameWikiAll := func(_ context.Context, cmd *cli.Command) error {
		logger.Info("started")

		client := cmd.Metadata["client"].(*wiki.Client)
		pages, err := client.List(cmd.String(projectKey.Name), cmd.String(pattern.Name))
		if err != nil {
			return err
		}

		for _, page := range pages {
			if err := client.Rename(page, cmd.String(oldString.Name), cmd.String(newString.Name)); err != nil {
				return err
			}
		}

		logger.Info("stopped")
		return nil
	}

	replaceWikiAll := func(_ context.Context, cmd *cli.Command) error {
		logger.Info("started")

		client := cmd.Metadata["client"].(*wiki.Client)
		pages, err := client.List(cmd.String(projectKey.Name), cmd.String(pattern.Name))
		if err != nil {
			return err
		}

		for _, page := range pages {
			detail, err := client.Get(page.ID)
			if err != nil {
				return err
			}
			if err := client.Replace(detail, cmd.StringSlice(pairs.Name)...); err != nil {
				return err
			}
		}

		logger.Info("stopped")
		return nil
	}

	return &cli.Command{
		Name:                  name,
		Version:               version.Version(),
		Usage:                 "Backlog utilities",
		Description:           "A cli application for Backlog utilities.",
		HideHelpCommand:       true,
		EnableShellCompletion: true,
		Writer:                w,
		ErrWriter:             ew,
		Commands: []*cli.Command{
			{
				Name:  "wiki",
				Usage: "Backlog wiki utilities",
				Commands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List wiki pages with optional pattern",
						Before: beforeWiki,
						Action: listWiki,
						Flags:  []cli.Flag{loglevel, baseURL, apiKey, projectKey, pattern},
					},
					{
						Name:   "rename",
						Usage:  "Rename wiki page",
						Before: beforeWiki,
						Action: renameWiki,
						Flags:  []cli.Flag{loglevel, baseURL, apiKey, wikiID, oldString, newString},
					},
					{
						Name:   "replace",
						Usage:  "Replace strings in the content of wiki page",
						Before: beforeWiki,
						Action: replaceWiki,
						Flags:  []cli.Flag{loglevel, baseURL, apiKey, wikiID, pairs},
					},
					{
						Name:   "rename-all",
						Usage:  "List wiki pages and rename them with optional pattern",
						Before: beforeWiki,
						Action: renameWikiAll,
						Flags:  []cli.Flag{loglevel, baseURL, apiKey, projectKey, pattern, oldString, newString},
					},
					{
						Name:   "replace-all",
						Usage:  "List wiki pages and replace strings in the content with optional pattern",
						Before: beforeWiki,
						Action: replaceWikiAll,
						Flags:  []cli.Flag{loglevel, baseURL, apiKey, projectKey, pattern, pairs},
					},
				},
			},
		},
	}
}
