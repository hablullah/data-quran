package qurancom

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

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

	// Process words
	err = processWords(ctx, cacheDir, dstDir)
	if err != nil {
		return err
	}

	return nil
}

func cleanDstDir(dstDir string) error {
	return filepath.WalkDir(dstDir, func(path string, d fs.DirEntry, err error) error {
		// Remove all file suffixed with "-qurancom.json"
		dName := d.Name()
		if d.IsDir() || !strings.HasSuffix(dName, "-qurancom.json") {
			return nil
		}

		dDir := filepath.Base(filepath.Dir(path))
		switch dDir {
		case "surah-info",
			"surah-translation",
			"word-text",
			"word-translation",
			"word-transliteration":
			return os.Remove(path)
		}

		return nil
	})
}
