package tanzilText

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "tanzil-text",
		Action: cliAction,
		Usage:  "download Quran text from Tanzil.net",
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
	cacheDir := filepath.Join(dstDir, ".cache", "tanzil-text")
	if c.Bool("clear-cache") {
		os.RemoveAll(cacheDir)
	}
	os.MkdirAll(cacheDir, os.ModePerm)

	// Create download URLs
	quranURLs := createQuranURLs()

	// Filter download request that not cached
	var requests []dl.Request
	for _, r := range quranURLs {
		dstPath := filepath.Join(cacheDir, r.FileName)
		if !util.FileExist(dstPath) {
			requests = append(requests, r)
		}
	}

	// Batch download the request
	ctx := context.Background()
	err := dl.BatchDownload(ctx, cacheDir, requests)
	if err != nil {
		return err
	}

	// Parse all file
	dataList, err := parse(cacheDir)
	if err != nil {
		return err
	}

	// Write to file
	err = cleanDstDir(dstDir)
	if err != nil {
		return err
	}

	err = writeTexts(dstDir, dataList)
	if err != nil {
		return err
	}

	return nil
}
