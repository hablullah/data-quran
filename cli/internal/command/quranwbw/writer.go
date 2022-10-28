package quranwbw

import (
	"data-quran-cli/internal/util"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func writeData(dstDir string, arabicDataList []ArabicData) error {
	logrus.Printf("writing data")

	// Prepare destination path
	dstDir = filepath.Join(dstDir, "word")
	os.MkdirAll(dstDir, os.ModePerm)

	// Prepare data
	type wordData struct {
		Surah    int `json:"surah"`
		Ayah     int `json:"ayah"`
		Position int `json:"position"`
	}

	output := map[string]wordData{}
	for i, data := range arabicDataList {
		key := fmt.Sprintf("%05d", i+1)
		output[key] = wordData{
			Surah:    data.Surah,
			Ayah:     data.Ayah,
			Position: data.Position,
		}
	}

	// Encode data to file
	dstPath := filepath.Join(dstDir, "word.json")
	err := util.EncodeSortedKeyJson(dstPath, &output)
	if err != nil {
		return fmt.Errorf("create file for word failed: %w", err)
	}

	return nil
}

func writeTexts(dstDir string, arabicDataList []ArabicData, textType string) error {
	logrus.Printf("writing text for %s", textType)

	// Prepare destination dir and path
	dstBaseDir := "word-text"
	dstPath := fmt.Sprintf("%s-quranwbw.json", textType)
	if textType == "transliteration" {
		dstBaseDir = "word-transliteration"
		dstPath = "en-quranwbw.json"
	}

	// Create dst dir
	dstDir = filepath.Join(dstDir, dstBaseDir)
	os.MkdirAll(dstDir, os.ModePerm)

	// Prepare data
	output := map[string]string{}
	for i, data := range arabicDataList {
		var text string
		switch textType {
		case "nastaliq":
			text = data.Nastaliq
		case "uthmani":
			text = data.Uthmani
		case "transliteration":
			text = data.Transliteration
		}

		key := fmt.Sprintf("%05d", i+1)
		output[key] = text
	}

	// Encode data to file
	dstPath = filepath.Join(dstDir, dstPath)
	err := util.EncodeSortedKeyJson(dstPath, &output)
	if err != nil {
		return fmt.Errorf("create file for text %s failed: %w", textType, err)
	}

	return nil
}

func writeTranslations(dstDir string, language, languageID string, translations []string) error {
	logrus.Printf("writing translation for %s", language)

	// Prepare destination dir
	dstDir = filepath.Join(dstDir, "word-translation")
	os.MkdirAll(dstDir, os.ModePerm)

	// Prepare destination path
	dstPath := fmt.Sprintf("%s-quranwbw.json", languageID)
	dstPath = filepath.Join(dstDir, dstPath)

	// Prepare output data
	output := map[string]string{}
	for i, trans := range translations {
		key := fmt.Sprintf("%05d", i+1)
		output[key] = trans
	}

	// Encode data to file
	dstPath = filepath.Join(dstDir, dstPath)
	err := util.EncodeSortedKeyJson(dstPath, &output)
	if err != nil {
		return fmt.Errorf("create file for trans %s failed: %w", language, err)
	}

	return nil
}
