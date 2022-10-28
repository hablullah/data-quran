package quranwbw

import (
	"data-quran-cli/internal/norm"
	"data-quran-cli/internal/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ArabicData struct {
	Surah           int
	Ayah            int
	Position        int
	Nastaliq        string
	Uthmani         string
	Transliteration string
}

func parseAllArabic(cacheDir string) ([]ArabicData, map[string]int, error) {
	var allData []ArabicData
	allWordCounts := map[string]int{}

	for surah := 1; surah <= 114; surah++ {
		surahPath := fmt.Sprintf("arabic-%03d.json", surah)
		surahPath = filepath.Join(cacheDir, surahPath)

		dataList, wordCounts, err := parseArabic(surahPath, surah)
		if err != nil {
			return nil, nil, err
		}

		allData = append(allData, dataList...)
		for id, count := range wordCounts {
			allWordCounts[id] = count
		}
	}

	// Make sure there are 77,429 words
	if n := len(allData); n != nWords {
		err := fmt.Errorf("wrong count of words, expected %d got %d", nWords, n)
		return nil, nil, err
	}

	return allData, allWordCounts, nil
}

func parseArabic(srcPath string, surah int) ([]ArabicData, map[string]int, error) {
	// Open file
	srcName := filepath.Base(srcPath)
	src, err := os.Open(srcPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open %s: %w", srcName, err)
	}
	defer src.Close()

	// Decode source file
	var srcData map[int]string
	err = json.NewDecoder(src).Decode(&srcData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode %s: %w", srcName, err)
	}

	// Process and normalize data
	var dataList []ArabicData
	wordCounts := map[string]int{}
	nAyah := util.ListSurah[surah].NAyah

	for ayah := 1; ayah <= nAyah; ayah++ {
		str := srcData[ayah]
		str = norm.NormalizeUnicode(str)
		words := strings.Split(str, "//")

		for pos, word := range words {
			parts := strings.Split(word, "/")
			if len(parts) != 3 {
				err = fmt.Errorf("arabic word in %d:%d:%d doesn't has three parts", surah, ayah, pos)
				return nil, nil, err
			}

			dataList = append(dataList, ArabicData{
				Surah:           surah,
				Ayah:            ayah,
				Position:        pos,
				Nastaliq:        strings.TrimSpace(parts[0]),
				Uthmani:         strings.TrimSpace(parts[1]),
				Transliteration: strings.TrimSpace(parts[2]),
			})
		}

		ayahID := fmt.Sprintf("%d-%d", surah, ayah)
		wordCounts[ayahID] = len(words)
	}

	return dataList, wordCounts, nil
}
