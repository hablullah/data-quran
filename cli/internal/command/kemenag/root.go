package kemenag

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"os"
	"path/filepath"
	"time"

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
	opts := dl.BatchOption{NWorker: 1, Delay: time.Second}
	err := dl.BatchDownload(ctx, cacheDir, requests, &opts)
	if err != nil {
		return err
	}

	// Clean dst dir
	if err = cleanDstDir(dstDir); err != nil {
		return err
	}

	// Parse and write list surah translation
	err = parseAndWriteListSurah(cacheDir, dstDir)
	if err != nil {
		return err
	}

	// Parse surah and write all ayah translation
	err = parseAndWriteAllSurah(cacheDir, dstDir)
	if err != nil {
		return err
	}

	// Parse all tafsir
	tafsirs, err := parseAllTafsir(cacheDir)
	if err != nil {
		return err
	}

	// Split tahlili and wajiz tafsir
	wajizTafsirs := make([]string, len(tafsirs))
	tahliliTafsirs := make([]string, len(tafsirs))
	for i, t := range tafsirs {
		wajizTafsirs[i] = t.TafsirWajiz
		tahliliTafsirs[i] = t.TafsirTahlili
	}

	// Write tafsir
	tafsirNames := map[string][]string{
		"id-ringkas-kemenag": wajizTafsirs,

		// TODO: for now we don't generate tafsir tahlili because
		// it still has a lot of typos and weird unicode errors.
		// "id-tahlili-kemenag": tahliliTafsirs,
	}

	for name, tafsirs := range tafsirNames {
		err = writeTafsir(dstDir, name, tafsirs)
		if err != nil {
			return err
		}
	}

	return nil
}
