package quranwbw

import (
	"data-quran-cli/internal/util"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

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

func writeData(dstDir string, arabicDataList []ArabicOutput) error {
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

func writeTexts(dstDir string, arabicDataList []ArabicOutput, textType string) error {
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

func writeOverviews(dstDir string, overviews map[int]string) error {
	logrus.Println("writing surah overviews")

	// Prepare destination dir
	dstDir = filepath.Join(dstDir, "surah-info")
	os.MkdirAll(dstDir, os.ModePerm)

	// Open destination file
	dstPath := filepath.Join(dstDir, "en-quranwbw.md")
	f, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create file for surah overview: %w", err)
	}
	defer f.Close()

	// Write metadata
	f.WriteString("<!--\n")
	f.WriteString("Language: english\n")
	f.WriteString("Source  : QuranWbW.com\n")
	f.WriteString("-->\n\n")

	// Write each info
	for surah := 1; surah <= 114; surah++ {
		f.WriteString("# ")
		f.WriteString(strconv.Itoa(surah))
		f.WriteString("\n\n")

		text := overviews[surah]
		if text == "" {
			f.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		}

		f.WriteString(text)
		f.WriteString("\n\n")
	}

	// Flush the data
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("write file for surah overview failed: %w", err)
	}

	return nil
}
