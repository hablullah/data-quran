package quranwbw

import (
	"data-quran-cli/internal/util"
	"fmt"
	"path/filepath"
	"strings"
)

type ArabicInput struct {
	P int    `json:"p"`
	W string `json:"w"`
	E string `json:"e"`
}

type ArabicOutput struct {
	Surah           int
	Ayah            int
	Position        int
	Nastaliq        string
	Uthmani         string
	Transliteration string
}

func parseAllArabic(cacheDir string) ([]ArabicOutput, map[string]int, error) {
	var allData []ArabicOutput
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

func parseArabic(srcPath string, surah int) ([]ArabicOutput, map[string]int, error) {
	// Decode source file
	var srcData map[int]ArabicInput
	srcName := filepath.Base(srcPath)
	err := util.DecodeJsonFile(srcPath, &srcData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode %s: %w", srcName, err)
	}

	// Process and normalize data
	var dataList []ArabicOutput
	wordCounts := map[string]int{}
	nAyah := util.ListSurah[surah].NAyah

	for ayah := 1; ayah <= nAyah; ayah++ {
		str := srcData[ayah].W
		words := strings.Split(str, "|")

		for pos, word := range words {
			parts := strings.Split(word, "/")
			if len(parts) != 4 {
				err = fmt.Errorf("arabic word in %d:%d:%d doesn't has 4 parts", surah, ayah, pos)
				return nil, nil, err
			}

			dataList = append(dataList, ArabicOutput{
				Surah:           surah,
				Ayah:            ayah,
				Position:        pos,
				Nastaliq:        strings.TrimSpace(parts[0]),
				Uthmani:         strings.TrimSpace(parts[1]),
				Transliteration: strings.TrimSpace(parts[3]),
			})
		}

		ayahID := fmt.Sprintf("%d-%d", surah, ayah)
		wordCounts[ayahID] = len(words)
	}

	return dataList, wordCounts, nil
}
