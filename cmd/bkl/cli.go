package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/nekrassov01/backlog-utils/backlog"
	"github.com/nekrassov01/backlog-utils/backlog/wiki"
	"github.com/urfave/cli/v3"
)

const name = "bkl"

var logger = &slog.Logger{}

type app struct {
	*cli.Command
	loglevel   *cli.StringFlag
	baseURL    *cli.StringFlag
	apiKey     *cli.StringFlag
	projectKey *cli.StringFlag
	pattern    *cli.StringFlag
	wikiID     *cli.IntFlag
	old        *cli.StringFlag
	new        *cli.StringFlag
	pairs      *cli.StringSliceFlag
}

func newApp(w, ew io.Writer) *app {
	logger = backlog.NewLogger(w, slog.LevelInfo.String())

	a := app{}
	a.loglevel = &cli.StringFlag{
		Name:    "log-level",
		Usage:   "set log level",
		Sources: cli.EnvVars("BACKLOG_LOG_LEVEL"),
		Value:   slog.LevelInfo.String(),
	}
	a.baseURL = &cli.StringFlag{
		Name:    "base-url",
		Usage:   "set backlog base url",
		Sources: cli.EnvVars("BACKLOG_URL"),
	}
	a.apiKey = &cli.StringFlag{
		Name:    "api-key",
		Usage:   "set backlog api key",
		Sources: cli.EnvVars("BACKLOG_API_KEY"),
	}
	a.projectKey = &cli.StringFlag{
		Name:     "project-key",
		Usage:    "set backlog project key",
		Required: true,
	}
	a.pattern = &cli.StringFlag{
		Name:  "pattern",
		Usage: "set pattern to search for wiki pages",
	}
	a.wikiID = &cli.IntFlag{
		Name:     "wiki-id",
		Usage:    "set backlog wiki id",
		Required: true,
	}
	a.old = &cli.StringFlag{
		Name:     "old",
		Usage:    "set string to be replaced in wiki page",
		Required: true,
	}
	a.new = &cli.StringFlag{
		Name:     "new",
		Usage:    "set new string after replacement in wiki page",
		Required: true,
	}
	a.pairs = &cli.StringSliceFlag{
		Name:     "pairs",
		Usage:    "set pairs of old and new repalacements for wiki page",
		Required: true,
	}

	a.Command = &cli.Command{
		Name:                  name,
		Version:               getVersion(),
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
						Before: a.before,
						Action: a.listWiki,
						Flags:  []cli.Flag{a.loglevel, a.baseURL, a.apiKey, a.projectKey, a.pattern},
					},
					{
						Name:   "rename",
						Usage:  "Rename wiki page",
						Before: a.before,
						Action: a.renameWiki,
						Flags:  []cli.Flag{a.loglevel, a.baseURL, a.apiKey, a.wikiID, a.old, a.new},
					},
					{
						Name:   "replace",
						Usage:  "Replace strings in the content of wiki page",
						Before: a.before,
						Action: a.replaceWiki,
						Flags:  []cli.Flag{a.loglevel, a.baseURL, a.apiKey, a.wikiID, a.pairs},
					},
					{
						Name:   "rename-all",
						Usage:  "List wiki pages and rename them with optional pattern",
						Before: a.before,
						Action: a.renameWikiAll,
						Flags:  []cli.Flag{a.loglevel, a.baseURL, a.apiKey, a.projectKey, a.pattern, a.old, a.new},
					},
					{
						Name:   "replace-all",
						Usage:  "List wiki pages and replace strings in the content with optional pattern",
						Before: a.before,
						Action: a.replaceWikiAll,
						Flags:  []cli.Flag{a.loglevel, a.baseURL, a.apiKey, a.projectKey, a.pattern, a.pairs},
					},
				},
			},
		},
	}

	return &a
}

func (a *app) before(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	level := cmd.String(a.loglevel.Name)
	if level == slog.LevelInfo.String() {
		return ctx, nil
	}
	logger = backlog.NewLogger(a.Writer, level)
	return ctx, nil
}

func (a *app) listWiki(_ context.Context, cmd *cli.Command) error {
	logger.Info("started")

	url := cmd.String(a.baseURL.Name)
	apiKey := cmd.String(a.apiKey.Name)
	projectKey := cmd.String(a.projectKey.Name)
	pattern := cmd.String(a.pattern.Name)

	client, err := wiki.New(a.Writer, url, apiKey)
	if err != nil {
		return err
	}

	pages, err := client.List(projectKey, pattern)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(a.Writer)
	for _, page := range pages {
		if err := enc.Encode(page); err != nil {
			return err
		}
	}

	logger.Info("stopped")
	return nil
}

func (a *app) renameWiki(_ context.Context, cmd *cli.Command) error {
	logger.Info("started")

	url := cmd.String(a.baseURL.Name)
	apiKey := cmd.String(a.apiKey.Name)
	wikiID := cmd.Int(a.wikiID.Name)
	before := cmd.String(a.old.Name)
	after := cmd.String(a.new.Name)

	client, err := wiki.New(a.Writer, url, apiKey)
	if err != nil {
		return err
	}

	page, err := client.Get(wikiID)
	if err != nil {
		return err
	}

	if err := client.Rename(page, before, after); err != nil {
		return err
	}

	logger.Info("stopped")
	return nil
}

func (a *app) renameWikiAll(_ context.Context, cmd *cli.Command) error {
	logger.Info("started")

	url := cmd.String(a.baseURL.Name)
	apiKey := cmd.String(a.apiKey.Name)
	projectKey := cmd.String(a.projectKey.Name)
	pattern := cmd.String(a.pattern.Name)
	before := cmd.String(a.old.Name)
	after := cmd.String(a.new.Name)

	client, err := wiki.New(a.Writer, url, apiKey)
	if err != nil {
		return err
	}

	pages, err := client.List(projectKey, pattern)
	if err != nil {
		return err
	}

	for _, page := range pages {
		if err := client.Rename(page, before, after); err != nil {
			return err
		}
	}

	logger.Info("stopped")
	return nil
}

func (a *app) replaceWiki(_ context.Context, cmd *cli.Command) error {
	logger.Info("started")

	url := cmd.String(a.baseURL.Name)
	apiKey := cmd.String(a.apiKey.Name)
	wikiID := cmd.Int(a.wikiID.Name)
	pairs := cmd.StringSlice(a.pairs.Name)

	client, err := wiki.New(a.Writer, url, apiKey)
	if err != nil {
		return err
	}

	page, err := client.Get(wikiID) // Get page content
	if err != nil {
		return err
	}

	if err := client.Replace(page, pairs...); err != nil {
		return err
	}

	logger.Info("stopped")
	return nil
}

func (a *app) replaceWikiAll(_ context.Context, cmd *cli.Command) error {
	logger.Info("started")

	url := cmd.String(a.baseURL.Name)
	apiKey := cmd.String(a.apiKey.Name)
	projectKey := cmd.String(a.projectKey.Name)
	pattern := cmd.String(a.pattern.Name)
	pairs := cmd.StringSlice(a.pairs.Name)

	client, err := wiki.New(a.Writer, url, apiKey)
	if err != nil {
		return err
	}

	pages, err := client.List(projectKey, pattern)
	if err != nil {
		return err
	}

	for _, page := range pages {
		detail, err := client.Get(page.ID) // Get page content
		if err != nil {
			return err
		}
		if err := client.Replace(detail, pairs...); err != nil {
			return err
		}
	}

	logger.Info("stopped")
	return nil
}
