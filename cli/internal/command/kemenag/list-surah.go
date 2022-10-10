package kemenag

import (
	"data-quran-cli/internal/norm"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type ListSurahContainer struct {
	Msg  string `json:"msg"`
	Data []struct {
		ID              int    `json:"id"`
		SuratName       string `json:"surat_name"`
		SuratText       string `json:"surat_text"`
		SuratTerjemahan string `json:"surat_terjemahan"`
		GolonganSurah   string `json:"golongan_surah"`
		CountAyat       int    `json:"count_ayat"`
	} `json:"data"`
}

type SurahOutput struct {
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

func parseAndWriteListSurah(cacheDir, dstDir string) error {
	logrus.Println("parse and write list surah")

	// Open and decode source file
	srcPath := filepath.Join(cacheDir, "list-surah.json")
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open list-surah: %w", err)
	}
	defer src.Close()

	var srcData ListSurahContainer
	err = json.NewDecoder(src).Decode(&srcData)
	if err != nil {
		return fmt.Errorf("failed to decode list-surah: %w", err)
	}

	// Normalize and convert the data
	outputs := make([]SurahOutput, len(srcData.Data))
	for i, d := range srcData.Data {
		outputs[i] = SurahOutput{
			Name:        norm.NormalizeUnicode(d.SuratName),
			Translation: norm.NormalizeUnicode(d.SuratTerjemahan),
		}
	}

	// Save as json
	dstDir = filepath.Join(dstDir, "surah-translation")
	if err = os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create list-surah dir: %w", err)
	}

	dstPath := filepath.Join(dstDir, "id-kemenag.json")
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create list-surah dst: %w", err)
	}
	defer dst.Close()

	encoder := json.NewEncoder(dst)
	encoder.SetIndent("", "    ")
	if err = encoder.Encode(&outputs); err != nil {
		return fmt.Errorf("failed to write list-surah: %w", err)
	}

	return nil
}
