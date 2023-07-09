package qurancom

import (
	"context"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var nWords = 77_429

func Command() *cli.Command {
	return &cli.Command{
		Name:   "qurancom",
		Action: cliAction,
		Usage:  "download data from quran.com",
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
	cacheDir := filepath.Join(dstDir, ".cache", "qurancom")
	if c.Bool("clear-cache") {
		os.RemoveAll(cacheDir)
	}
	os.MkdirAll(cacheDir, os.ModePerm)

	// Clean dst dir
	if err := cleanDstDir(dstDir); err != nil {
		return err
	}

	// Process chapter info
	ctx := context.Background()
	err := processChapterInfo(ctx, cacheDir, dstDir)
	if err != nil {
		return err
	}

	// Process chapter names
	err = processChapterList(ctx, cacheDir, dstDir)
	if err != nil {
		return err
	}

	// Download words
	err = downloadAllWords(ctx, cacheDir)
	if err != nil {
		return err
	}

	err = parseAndWriteWordText(cacheDir, dstDir)
	if err != nil {
		return err
	}

	err = parseAndWriteWordTransliteration(cacheDir, dstDir)
	if err != nil {
		return err
	}

	return nil
}
