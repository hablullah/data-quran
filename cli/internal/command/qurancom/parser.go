package qurancom

import (
	"data-quran-cli/internal/norm"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type WordSource struct {
	Verses []struct {
		ID              int         `json:"id"`
		VerseNumber     int         `json:"verse_number"`
		ChapterID       int         `json:"chapter_id"`
		VerseKey        string      `json:"verse_key"`
		TextIndopak     string      `json:"text_indopak"`
		JuzNumber       int         `json:"juz_number"`
		HizbNumber      int         `json:"hizb_number"`
		RubElHizbNumber int         `json:"rub_el_hizb_number"`
		SajdahNumber    interface{} `json:"sajdah_number"`
		PageNumber      int         `json:"page_number"`
		Sajdah          interface{} `json:"sajdah"`
		TextMadani      string      `json:"text_madani"`
		Words           []struct {
			ID              int    `json:"id"`
			Position        int    `json:"position"`
			TextIndopak     string `json:"text_indopak"`
			VerseKey        string `json:"verse_key"`
			LineNumber      int    `json:"line_number"`
			PageNumber      int    `json:"page_number"`
			Code            string `json:"code"`
			ClassName       string `json:"class_name"`
			TextMadani      string `json:"text_madani"`
			CharType        string `json:"char_type"`
			Transliteration struct {
				Text         string `json:"text"`
				LanguageName string `json:"language_name"`
			} `json:"transliteration"`
			Translation struct {
				LanguageName string `json:"language_name"`
				Text         string `json:"text"`
			} `json:"translation"`
			Audio struct {
				URL string `json:"url"`
			} `json:"audio"`
		} `json:"words"`
	} `json:"verses"`
	Pagination struct {
		CurrentPage int         `json:"current_page"`
		NextPage    interface{} `json:"next_page"`
		PrevPage    interface{} `json:"prev_page"`
		TotalPages  int         `json:"total_pages"`
		TotalCount  int         `json:"total_count"`
	} `json:"pagination"`
}

type WordText struct {
	Indopak         string
	Madani          string
	Transliteration string
}

func parseAllWordTranslations(cacheDir, language string) (map[string]string, error) {
	logrus.Printf("parsing word for %s", language)

	// Extract each surah in this language
	var id int
	translations := map[string]string{}

	for surah := 1; surah <= 114; surah++ {
		sTranslations, err := parseWordTranslations(cacheDir, language, surah)
		if err != nil {
			return nil, err
		}

		for i := range sTranslations {
			id++
			key := fmt.Sprintf("%05d", id)
			translations[key] = sTranslations[i]
		}
	}

	// Check if info complete
	if n := len(translations); n != nWords {
		logrus.Warnf("word trans for %s: want %d got %d", language, nWords, n)
		return nil, nil
	}

	return translations, nil
}

func parseWordTranslations(cacheDir, language string, surah int) ([]string, error) {
	// Open file
	path := fmt.Sprintf("word-%s-%03d.json", language, surah)
	path = filepath.Join(cacheDir, path)

	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("fail to open word trans for %s %d: %w", language, surah, err)
		return nil, err
	}
	defer f.Close()

	// Decode data
	var srcData WordSource
	err = json.NewDecoder(f).Decode(&srcData)
	if err != nil {
		err = fmt.Errorf("fail to decode word trans for %s %d: %w", language, surah, err)
		return nil, err
	}

	// Get list of word tranlations
	var translations []string
	for _, verse := range srcData.Verses {
		for _, word := range verse.Words {
			// We only care about word
			if word.CharType != "word" {
				continue
			}

			// In Quran.com, if data for a language not exist, they will
			// fallback into using English language. In this case, we
			// will just put missing mark.
			if word.Translation.LanguageName == "english" && language != "en" {
				translations = append(translations, "[[MISSING]]")
			} else {
				// Normalize and clean
				wordTrans := word.Translation.Text
				wordTrans = norm.NormalizeUnicode(wordTrans)
				if wordTrans == "" {
					wordTrans = "[[MISSING]]"
				}

				translations = append(translations, wordTrans)
			}
		}
	}

	return translations, nil
}

func parseAllWordTexts(cacheDir string) (map[string]WordText, error) {
	logrus.Printf("parsing word texts")

	// Extract each surah
	var id int
	texts := map[string]WordText{}

	for surah := 1; surah <= 114; surah++ {
		sTexts, err := parseWordTexts(cacheDir, surah)
		if err != nil {
			return nil, err
		}

		for _, sText := range sTexts {
			id++
			key := fmt.Sprintf("%05d", id)
			texts[key] = sText
		}
	}

	// Check if text complete
	if n := len(texts); n != nWords {
		logrus.Warnf("word text: want %d got %d", nWords, n)
		return nil, nil
	}

	return texts, nil
}

func parseWordTexts(cacheDir string, surah int) ([]WordText, error) {
	// Open Englishfile
	path := fmt.Sprintf("word-en-%03d.json", surah)
	path = filepath.Join(cacheDir, path)

	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("fail to open word text for %d: %w", surah, err)
		return nil, err
	}
	defer f.Close()

	// Decode data
	var srcData WordSource
	err = json.NewDecoder(f).Decode(&srcData)
	if err != nil {
		err = fmt.Errorf("fail to decode word text for %d: %w", surah, err)
		return nil, err
	}

	// Get list of word texts
	var texts []WordText
	for _, verse := range srcData.Verses {
		for _, word := range verse.Words {
			// We only care about word
			if word.CharType != "word" {
				continue
			}

			texts = append(texts, WordText{
				Madani:          norm.NormalizeUnicode(word.TextMadani),
				Indopak:         norm.NormalizeUnicode(word.TextIndopak),
				Transliteration: norm.NormalizeUnicode(word.Transliteration.Text),
			})
		}
	}

	return texts, nil
}
