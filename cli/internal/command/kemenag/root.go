package kemenag

import (
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

	// Download data
	err := downloadListSurah(cacheDir)
	if err != nil {
		return err
	}

	err = downladAllTafsir(cacheDir)
	if err != nil {
		return err
	}

	// Parse and write list surah
	err = parseAndWriteListSurah(cacheDir, dstDir)
	if err != nil {
		return err
	}

	// Parse each surah to extract text, trans and tafsirs
	listAyah, err := parseAllSurah(cacheDir)
	if err != nil {
		return err
	}

	err = writeQuranBasicData(dstDir, listAyah, TextArabic)
	if err != nil {
		return err
	}

	err = writeQuranBasicData(dstDir, listAyah, Transliteration)
	if err != nil {
		return err
	}

	err = writeQuranBasicData(dstDir, listAyah, TafsirWajiz)
	if err != nil {
		return err
	}

	err = writeQuranBasicData(dstDir, listAyah, TafsirTahlili)
	if err != nil {
		return err
	}

	err = writeQuranTranslation(dstDir, listAyah)
	if err != nil {
		return err
	}

	err = writeSurahInfo(dstDir, listAyah)
	if err != nil {
		return err
	}

	return nil
}
