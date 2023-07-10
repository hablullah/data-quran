package kemenag

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	urlListSurah = "https://web-api.qurankemenag.net/quran-surah"
)

func processListSurah(ctx context.Context, cacheDir, dstDir string) error {
	err := downloadListSurah(ctx, cacheDir)
	if err != nil {
		return err
	}

	listSurah, err := parseListSurah(cacheDir)
	if err != nil {
		return err
	}

	err = writeListSurah(dstDir, listSurah)
	if err != nil {
		return err
	}

	return nil
}

func downloadListSurah(ctx context.Context, cacheDir string) error {
	logrus.Printf("downloading list surah")
	dstPath := filepath.Join(cacheDir, "list-surah.json")

	if !util.FileExist(dstPath) {
		req := dl.Request{URL: urlListSurah}
		err := dl.Download(ctx, http.DefaultClient, dstPath, req)
		if err != nil {
			return fmt.Errorf("failed to download list surah: %w", err)
		}
	}

	return nil
}

func parseListSurah(cacheDir string) ([]Surah, error) {
	logrus.Printf("parsing list surah")
	var listSurah struct {
		Data []Surah `json:"data"`
	}

	srcPath := filepath.Join(cacheDir, "list-surah.json")
	err := util.DecodeJsonFile(srcPath, &listSurah)
	if err != nil {
		return nil, fmt.Errorf("failed to decode list surah: %w", err)
	}

	return listSurah.Data, nil
}

func writeListSurah(dstDir string, listSurah []Surah) error {
	logrus.Printf("writing list surah")

	// Prepare data
	data := make(map[string]ListSurahEntry)
	for i, s := range listSurah {
		ayahId := fmt.Sprintf("%04d", i+1)
		data[ayahId] = ListSurahEntry{
			Name:        strings.TrimSpace(s.Transliteration),
			Translation: strings.TrimSpace(s.Translation),
		}
	}

	// Write to file
	dstPath := filepath.Join(dstDir, "surah-translation", "id-kemenag.json")
	err := util.EncodeSortedKeyJson(dstPath, &data)
	if err != nil {
		return fmt.Errorf("failed to write list surah: %w", err)
	}

	return nil
}
