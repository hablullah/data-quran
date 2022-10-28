package quranwbw

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var nWords = 77_429

func Command() *cli.Command {
	return &cli.Command{
		Name:   "quranwbw",
		Action: cliAction,
		Usage:  "download data from quranwbw",
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
	cacheDir := filepath.Join(dstDir, ".cache", "quranwbw")
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

	// Parse Arabic data
	arabicDataList, wordCounts, err := parseAllArabic(cacheDir)
	if err != nil {
		return err
	}

	// Write word metadata, Arabic text and transliteration
	err = writeData(dstDir, arabicDataList)
	if err != nil {
		return err
	}

	err = writeTexts(dstDir, arabicDataList, "uthmani")
	if err != nil {
		return err
	}

	err = writeTexts(dstDir, arabicDataList, "nastaliq")
	if err != nil {
		return err
	}

	err = writeTexts(dstDir, arabicDataList, "transliteration")
	if err != nil {
		return err
	}

	// Parse and write translations
	for lang, langID := range languages {
		// Skip Arabic
		if lang == "arabic" {
			continue
		}

		// Parse translation
		translations, err := parseAllTranslationFiles(cacheDir, lang, wordCounts)
		if err != nil {
			return err
		}

		// Write translation
		err = writeTranslations(dstDir, lang, langID, translations)
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanDstDir(dstDir string) error {
	return filepath.WalkDir(dstDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		// Remove word.json
		dName := d.Name()
		dDir := filepath.Base(filepath.Dir(path))
		if dDir == "word" && dName == "word.json" {
			return os.Remove(path)
		}

		// Remove all file suffixed with "-quranwbw.*"
		switch dDir {
		case "word-text",
			"word-translation",
			"word-transliteration":
			if strings.HasSuffix(dName, "-quranwbw.json") {
				return os.Remove(path)
			}
		}

		return nil
	})
}
