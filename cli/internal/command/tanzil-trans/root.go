package tanzilTrans

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "tanzil-trans",
		Action: cliAction,
		Usage:  "download translation data from Tanzil.net",
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
	cacheDir := filepath.Join(dstDir, ".cache", "tanzil-trans")
	if c.Bool("clear-cache") {
		os.RemoveAll(cacheDir)
	}
	os.MkdirAll(cacheDir, os.ModePerm)

	// Download translation page from Tanzil
	err := downloadTranslationPage(cacheDir)
	if err != nil {
		return err
	}

	// Parse URLs from the translation page
	downloadRequests, err := parseTranslationPage(cacheDir)
	if err != nil {
		return err
	}

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
	err = dl.BatchDownload(ctx, cacheDir, requests, nil)
	if err != nil {
		return err
	}

	// Parse all file
	dataList, err := parse(cacheDir)
	if err != nil {
		return err
	}

	// Clean data
	for i, data := range dataList {
		if fnCleaner, exist := cleanerList[data.FileName]; exist {
			logrus.Printf("cleaning %s", data.FileName)
			dataList[i] = fnCleaner(data)
		}
	}

	// Write to file
	err = cleanDstDir(dstDir)
	if err != nil {
		return err
	}

	err = writeTranslations(dstDir, dataList)
	if err != nil {
		return err
	}

	return nil
}
