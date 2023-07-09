package qurancom

import (
	"data-quran-cli/internal/util"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func cleanDstDir(dstDir string) error {
	return filepath.WalkDir(dstDir, func(path string, d fs.DirEntry, err error) error {
		// Remove all file suffixed with "-qurancom.json"
		dName := d.Name()
		if d.IsDir() || !strings.HasSuffix(dName, "-qurancom.json") {
			return nil
		}

		dDir := filepath.Base(filepath.Dir(path))
		switch dDir {
		case "surah-info",
			"surah-translation",
			"word-text",
			"word-translation",
			"word-transliteration":
			return os.Remove(path)
		}

		return nil
	})
}

func writeListSurah(dstDir string, language string, data map[string]ListSurahOutput) error {
	// If data is empty, stop
	if len(data) == 0 {
		return nil
	}

	logrus.Printf("writing surah list for %s", language)

	// Prepare destination path
	dstDir = filepath.Join(dstDir, "surah-translation")
	os.MkdirAll(dstDir, os.ModePerm)

	dstPath := fmt.Sprintf("%s-qurancom.json", language)
	dstPath = filepath.Join(dstDir, dstPath)

	// Encode data
	err := util.EncodeSortedKeyJson(dstPath, &data)
	if err != nil {
		return fmt.Errorf("fail to write surah list for %s: %w", language, err)
	}

	return nil
}

func writeWordTranslations(dstDir string, language string, translations map[string]string) error {
	// If data is empty, stop
	if len(translations) == 0 {
		return nil
	}

	logrus.Printf("writing word translation for %s", language)
	dstPath := fmt.Sprintf("%s-qurancom.json", language)
	dstPath = filepath.Join(dstDir, "word-translation", dstPath)
	err := writeWordJson(dstPath, &translations)
	if err != nil {
		return fmt.Errorf("create file for word trans %s failed: %w", language, err)
	}

	return nil
}

func writeWordTexts(dstDir string, texts map[string]WordText) error {
	// If data is empty, stop
	if len(texts) == 0 {
		return nil
	}

	// Split data
	madaniTexts := map[string]string{}
	indopakTexts := map[string]string{}
	transliterations := map[string]string{}

	for key, text := range texts {
		madaniTexts[key] = text.Madani
		indopakTexts[key] = text.Indopak
		transliterations[key] = text.Transliteration
	}

	// Write Madani
	logrus.Printf("writing word texts Madani")
	dstPath := filepath.Join(dstDir, "word-text", "madani-qurancom.json")
	err := writeWordJson(dstPath, &madaniTexts)
	if err != nil {
		return fmt.Errorf("create file for Madani text failed: %w", err)
	}

	// Write Indopak
	logrus.Printf("writing word texts Indopak")
	dstPath = filepath.Join(dstDir, "word-text", "indopak-qurancom.json")
	err = writeWordJson(dstPath, &indopakTexts)
	if err != nil {
		return fmt.Errorf("create file for Indopak text failed: %w", err)
	}

	// Write transliterations
	logrus.Printf("writing word transliterations")
	dstPath = filepath.Join(dstDir, "word-transliteration", "en-qurancom.json")
	err = writeWordJson(dstPath, &transliterations)
	if err != nil {
		return fmt.Errorf("create file for transliterations failed: %w", err)
	}

	return nil
}

func writeWordJson(path string, data any) error {
	// Prepare destination dir
	dstDir := filepath.Dir(path)
	os.MkdirAll(dstDir, os.ModePerm)

	// Encode data
	return util.EncodeSortedKeyJson(path, data)
}
