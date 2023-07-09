package qurancom

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	nurl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type WordResponse struct {
	Verses []struct {
		ID              int    `json:"id"`
		VerseNumber     int    `json:"verse_number"`
		VerseKey        string `json:"verse_key"`
		HizbNumber      int    `json:"hizb_number"`
		RubElHizbNumber int    `json:"rub_el_hizb_number"`
		RukuNumber      int    `json:"ruku_number"`
		ManzilNumber    int    `json:"manzil_number"`
		SajdahNumber    any    `json:"sajdah_number"`
		PageNumber      int    `json:"page_number"`
		JuzNumber       int    `json:"juz_number"`
		Words           []struct {
			ID           int    `json:"id"`
			Position     int    `json:"position"`
			AudioURL     string `json:"audio_url"`
			CharTypeName string `json:"char_type_name"`
			TextUthmani  string `json:"text_uthmani"`
			TextIndopak  string `json:"text_indopak"`
			TextImlaei   string `json:"text_imlaei"`
			PageNumber   int    `json:"page_number"`
			LineNumber   int    `json:"line_number"`
			Text         string `json:"text"`
			Translation  struct {
				Text         string `json:"text"`
				LanguageName string `json:"language_name"`
			} `json:"translation"`
			Transliteration struct {
				Text         string `json:"text"`
				LanguageName string `json:"language_name"`
			} `json:"transliteration"`
		} `json:"words"`
	} `json:"verses"`
	Pagination struct {
		PerPage      int `json:"per_page"`
		CurrentPage  int `json:"current_page"`
		NextPage     any `json:"next_page"`
		TotalPages   int `json:"total_pages"`
		TotalRecords int `json:"total_records"`
	} `json:"pagination"`
}

var (
	languagesForWord = map[string]string{
		"en":  "english",
		"ur":  "urdu",
		"id":  "indonesian",
		"bn":  "bengali",
		"tr":  "turkish",
		"fa":  "persian",
		"ru":  "russian",
		"hi":  "hindi",
		"de":  "german",
		"ta":  "tamil",
		"inh": "ingush",
	}
)

func urlWord(surah int, language string) string {
	query := nurl.Values{}
	query.Set("words", "1")
	query.Set("per_page", "300")
	query.Set("language", language)
	query.Set("word_fields", "text_uthmani,text_indopak,text_imlaei,translation,transliteration")

	rawURL := fmt.Sprintf("https://api.quran.com/api/v4/verses/by_chapter/%d", surah)
	url, _ := nurl.ParseRequestURI(rawURL)
	url.RawQuery = query.Encode()
	return url.String()
}

func processWords(ctx context.Context, cacheDir, dstDir string) error {
	// Download all words data
	err := downloadAllWords(ctx, cacheDir)
	if err != nil {
		return err
	}

	// Process word text
	err = parseAndWriteWordText(cacheDir, dstDir)
	if err != nil {
		return err
	}

	// Process word transliteration
	err = parseAndWriteWordTransliteration(cacheDir, dstDir)
	if err != nil {
		return err
	}

	// Process word translation
	err = parseAndWriteAllWordTranslation(cacheDir, dstDir)
	if err != nil {
		return err
	}

	return nil
}

func downloadWords(ctx context.Context, cacheDir string, language string) error {
	logrus.Printf("downloading words for %s", language)

	// Prepare download links
	var dlRequests []dl.Request
	for surah := 1; surah <= 114; surah++ {
		url := urlWord(surah, language)
		dstName := fmt.Sprintf("word-%s-%03d.json", language, surah)
		dstPath := filepath.Join(cacheDir, dstName)

		if !util.FileExist(dstPath) {
			dlRequests = append(dlRequests, dl.Request{FileName: dstName, URL: url})
		}
	}

	// Start batch download
	err := dl.BatchDownload(ctx, cacheDir, dlRequests, nil)
	if err != nil {
		return fmt.Errorf("failed to download words for %s: %w", language, err)
	}

	return nil
}

func downloadAllWords(ctx context.Context, cacheDir string) error {
	for lang := range languagesForWord {
		err := downloadWords(ctx, cacheDir, lang)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseAndWriteWordText(cacheDir, dstDir string) error {
	// Parse each surah
	var wordIdx int
	allUthmani := make(map[string]string)
	allIndopak := make(map[string]string)
	allImlaei := make(map[string]string)

	for surah := 1; surah <= 114; surah++ {
		srcPath := fmt.Sprintf("word-en-%03d.json", surah)
		srcPath = filepath.Join(cacheDir, srcPath)

		var src WordResponse
		err := util.DecodeJsonFile(srcPath, &src)
		if err != nil {
			return fmt.Errorf("failed to decode word text %d: %w", surah, err)
		}

		for _, verse := range src.Verses {
			for _, w := range verse.Words {
				if w.CharTypeName == "word" {
					wordIdx++
					strWordIdx := fmt.Sprintf("%05d", wordIdx)

					uthmani := strings.TrimSpace(w.TextUthmani)
					if uthmani == "" {
						uthmani = "[[MISSING]]"
					}

					indopak := strings.TrimSpace(w.TextIndopak)
					if indopak == "" {
						indopak = "[[MISSING]]"
					}

					imlaei := strings.TrimSpace(w.TextImlaei)
					if imlaei == "" {
						imlaei = "[[MISSING]]"
					}

					allUthmani[strWordIdx] = uthmani
					allIndopak[strWordIdx] = indopak
					allImlaei[strWordIdx] = imlaei
				}
			}
		}
	}

	// Make sure word count is 77429
	if wordIdx != 77429 {
		return fmt.Errorf("word text count 77529 != %d", wordIdx)
	}

	// Prepare directory
	logrus.Printf("writing word text")
	dstDir = filepath.Join(dstDir, "word-text")
	os.MkdirAll(dstDir, os.ModePerm)

	// Write uthmani
	dstPath := filepath.Join(dstDir, "uthmani-qurancom.json")
	err := util.EncodeSortedKeyJson(dstPath, &allUthmani)
	if err != nil {
		return fmt.Errorf("failed to write uthmani word: %w", err)
	}

	// Write indopak
	dstPath = filepath.Join(dstDir, "indopak-qurancom.json")
	err = util.EncodeSortedKeyJson(dstPath, &allIndopak)
	if err != nil {
		return fmt.Errorf("failed to write indopak word: %w", err)
	}

	// Write imlaei
	dstPath = filepath.Join(dstDir, "imlaei-qurancom.json")
	err = util.EncodeSortedKeyJson(dstPath, &allImlaei)
	if err != nil {
		return fmt.Errorf("failed to write imlaei word: %w", err)
	}

	return nil
}

func parseAndWriteWordTransliteration(cacheDir, dstDir string) error {
	// Parse each surah
	var wordIdx int
	allTransliteration := make(map[string]string)

	for surah := 1; surah <= 114; surah++ {
		srcPath := fmt.Sprintf("word-en-%03d.json", surah)
		srcPath = filepath.Join(cacheDir, srcPath)

		var src WordResponse
		err := util.DecodeJsonFile(srcPath, &src)
		if err != nil {
			return fmt.Errorf("failed to decode word transliteration %d: %w", surah, err)
		}

		for _, verse := range src.Verses {
			for _, w := range verse.Words {
				if w.CharTypeName == "word" {
					wordIdx++
					strWordIdx := fmt.Sprintf("%05d", wordIdx)
					translit := strings.TrimSpace(w.Transliteration.Text)
					if translit == "" {
						translit = "[[MISSING]]"
					}
					allTransliteration[strWordIdx] = translit
				}
			}
		}
	}

	// Make sure word count is 77429
	if wordIdx != 77429 {
		return fmt.Errorf("word transliteration count 77529 != %d", wordIdx)
	}

	// Prepare directory
	logrus.Printf("writing word transliteration")
	dstDir = filepath.Join(dstDir, "word-transliteration")
	os.MkdirAll(dstDir, os.ModePerm)

	// Write
	dstPath := filepath.Join(dstDir, "en-qurancom.json")
	err := util.EncodeSortedKeyJson(dstPath, &allTransliteration)
	if err != nil {
		return fmt.Errorf("failed to write word transliteration: %w", err)
	}

	return nil
}

func parseAndWriteWordTranslation(cacheDir, dstDir, language string) error {
	// Parse each surah
	var wordIdx int
	languageName := languagesForWord[language]
	allTranslation := make(map[string]string)

	for surah := 1; surah <= 114; surah++ {
		srcPath := fmt.Sprintf("word-%s-%03d.json", language, surah)
		srcPath = filepath.Join(cacheDir, srcPath)

		var src WordResponse
		err := util.DecodeJsonFile(srcPath, &src)
		if err != nil {
			return fmt.Errorf("failed to decode %q word trans %d: %w", language, surah, err)
		}

		for _, verse := range src.Verses {
			for _, w := range verse.Words {
				if w.CharTypeName == "word" {
					wordIdx++
					strWordIdx := fmt.Sprintf("%05d", wordIdx)
					trans := strings.TrimSpace(w.Translation.Text)
					if trans == "" || w.Translation.LanguageName != languageName {
						trans = "[[MISSING]]"
					}
					allTranslation[strWordIdx] = trans
				}
			}
		}
	}

	// Make sure word count is 77429
	if wordIdx != 77429 {
		return fmt.Errorf("word translation %q count 77529 != %d", language, wordIdx)
	}

	// Prepare directory
	logrus.Printf("writing word translation for %q", language)
	dstDir = filepath.Join(dstDir, "word-translation")
	os.MkdirAll(dstDir, os.ModePerm)

	// Write
	dstPath := fmt.Sprintf("%s-qurancom.json", language)
	dstPath = filepath.Join(dstDir, dstPath)
	err := util.EncodeSortedKeyJson(dstPath, &allTranslation)
	if err != nil {
		return fmt.Errorf("failed to write word translation %q: %w", language, err)
	}

	return nil
}

func parseAndWriteAllWordTranslation(cacheDir, dstDir string) error {
	for lang := range languagesForWord {
		err := parseAndWriteWordTranslation(cacheDir, dstDir, lang)
		if err != nil {
			return err
		}
	}
	return nil
}
