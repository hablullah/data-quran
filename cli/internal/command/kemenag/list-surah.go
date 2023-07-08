package kemenag

import (
	"data-quran-cli/internal/norm"
	"data-quran-cli/internal/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func parseAndWriteListSurah(cacheDir, dstDir string) error {
	listSurah, err := parseListSurah(cacheDir)
	if err != nil {
		return err
	}

	return writeListSurah(dstDir, listSurah)
}

func parseListSurah(cacheDir string) ([]Surah, error) {
	logrus.Printf("parsing list surah")

	// Open file
	srcPath := filepath.Join(cacheDir, "list-surah.json")
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read list surah: %w", err)
	}
	defer f.Close()

	// Decode JSON
	var listSurah struct {
		Data []Surah `json:"data"`
	}

	r := norm.NormalizeReader(f)
	err = json.NewDecoder(r).Decode(&listSurah)
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
