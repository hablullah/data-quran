package kemenag

import (
	"context"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "kemenag",
		Action: cliAction,
		Usage:  "download data from quran.kemenag.go.id",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dst",
				Aliases: []string{"d"},
				Usage:   "destination directory",
				Value:   ".",
			},
			&cli.BoolFlag{
				Name:    "clear-cache",
				Aliases: []string{"cc"},
				Usage:   "clear download cache",
			},
		},
	}
}

func cliAction(c *cli.Context) error {
	// Prepare cache dir
	dstDir := c.String("dst")
	cacheDir := filepath.Join(dstDir, ".cache", "kemenag")
	if c.Bool("clear-cache") {
		os.RemoveAll(cacheDir)
	}
	os.MkdirAll(cacheDir, os.ModePerm)

	// Process list surah
	ctx := context.Background()
	err := processListSurah(ctx, cacheDir, dstDir)
	if err != nil {
		return err
	}

	// Process tafsirs for each surah
	err = processTafsir(cacheDir, dstDir)
	if err != nil {
		return err
	}

	return nil
}
