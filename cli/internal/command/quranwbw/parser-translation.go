package quranwbw

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

func parseAllTranslationFiles(cacheDir string, language string, wordCounts map[string]int) ([]string, error) {
	var allTranslations []string

	for surah := 1; surah <= 114; surah++ {
		surahPath := fmt.Sprintf("%s-%03d.json", language, surah)
		surahPath = filepath.Join(cacheDir, surahPath)

		translations, err := parseTranslationFile(surahPath, surah, wordCounts)
		if err != nil {
			return nil, err
		}

		allTranslations = append(allTranslations, translations...)
	}

	// Make sure there are 77,429 words
	if n := len(allTranslations); n != nWords {
		return nil, fmt.Errorf("wrong count of trans for %s, expected %d got %d", language, nWords, n)
	}

	return allTranslations, nil
}

func parseTranslationFile(srcPath string, surah int, wordCounts map[string]int) ([]string, error) {
	// Open file
	srcName := filepath.Base(srcPath)
	src, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", srcName, err)
	}
	defer src.Close()

	// Decode source file
	var srcData map[int]string
	err = json.NewDecoder(src).Decode(&srcData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", srcName, err)
	}

	// Process and normalize data
	var translations []string
	nAyah := util.ListSurah[surah].NAyah

	for ayah := 1; ayah <= nAyah; ayah++ {
		str := srcData[ayah]
		str = norm.NormalizeUnicode(str)
		ayahTranslations := strings.Split(str, "//")

		// Check word count for this ayah
		ayahID := fmt.Sprintf("%d-%d", surah, ayah)
		nTranslations := len(ayahTranslations)
		expectedCount := wordCounts[ayahID]
		if nTranslations != expectedCount {
			logrus.Warnf("wrong word count for %s in %s: want %d got %d",
				srcName, ayahID, expectedCount, nTranslations)
		}

		// Save translations
		for i := 0; i < expectedCount; i++ {
			var wordTrans string
			if i < nTranslations {
				wordTrans = ayahTranslations[i]
				wordTrans = strings.TrimSpace(wordTrans)
			} else {
				wordTrans = "[[MISSING]]"
			}
			translations = append(translations, wordTrans)
		}
	}

	return translations, nil
}
