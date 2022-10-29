package islamhouse

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "islamhouse",
		Action: cliAction,
		Usage:  "download tafsir mokhtasar from IslamHouse.com",
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
	cacheDir := filepath.Join(dstDir, ".cache", "islamhouse")
	if c.Bool("clear-cache") {
		os.RemoveAll(cacheDir)
	}
	os.MkdirAll(cacheDir, os.ModePerm)

	// Create download URLs
	downloadRequests := createDownloadRequests()

	// Filter download request that not cached
	var requests []dl.Request
	for _, r := range downloadRequests {
		dstPath := filepath.Join(cacheDir, r.FileName)
		if !util.FileExist(dstPath) {
			requests = append(requests, r)
		}
	}

	// Batch download the request
	ctx := context.Background()
	err := dl.BatchDownload(ctx, cacheDir, requests, nil)
	if err != nil {
		return err
	}

	// Clean dst dir
	if err = cleanDstDir(dstDir); err != nil {
		return err
	}

	// Parse all tafsir page
	mapTafsirs := map[string][]string{}
	for _, src := range sourceNames {
		tafsirs, err := parseAllPages(cacheDir, src.Language)
		if err != nil {
			return err
		}
		mapTafsirs[src.Language] = tafsirs
	}

	// Write all tafsir page
	for _, src := range sourceNames {
		tafsirs := mapTafsirs[src.Language]
		name := fmt.Sprintf("%s-mokhtasar-islamhouse", src.Language)
		err = writeTafsirs(dstDir, name, tafsirs)
		if err != nil {
			return err
		}
	}

	return nil
}
