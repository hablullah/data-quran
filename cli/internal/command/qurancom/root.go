package qurancom

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"os"
	"path/filepath"

	md "github.com/JohannesKaufmann/html-to-markdown"
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

	// Parse list surah translation
	listSurahData := map[string]map[string]ListSurahOutput{}
	for _, lang := range languages {
		data, err := parseListSurah(cacheDir, lang)
		if err != nil {
			return err
		} else if len(data) > 0 {
			listSurahData[lang] = data
		}
	}

	// Write list surah translation
	for _, lang := range languages {
		data := listSurahData[lang]
		err = writeListSurah(dstDir, lang, data)
		if err != nil {
			return err
		}
	}

	// Parse all surah info
	mdc := md.NewConverter("", true, nil)
	listSurahInfo := map[string]*AllSurahInfoOutput{}
	for _, lang := range languages {
		data, err := parseAllSurahInfo(cacheDir, lang, mdc)
		if err != nil {
			return err
		} else if data != nil {
			listSurahInfo[lang] = data
		}
	}

	// Write surah info
	for _, lang := range languages {
		data := listSurahInfo[lang]
		err = writeSurahInfo(dstDir, lang, data)
		if err != nil {
			return err
		}
	}

	// Parse word translations
	listWordTranslations := map[string]map[string]string{}
	for _, lang := range languagesForWord {
		words, err := parseAllWords(cacheDir, lang)
		if err != nil {
			return err
		} else if len(words) > 0 {
			listWordTranslations[lang] = words
		}
	}

	// Write word translations
	for _, lang := range languagesForWord {
		translations := listWordTranslations[lang]
		err = writeWordTranslations(dstDir, lang, translations)
		if err != nil {
			return err
		}
	}

	return nil
}
