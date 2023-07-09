package qurancom

import (
	"data-quran-cli/internal/dl"
	"fmt"
)

var (
	listURL = "https://api.quran.com/api/v3/chapters?language=%s"
	infoURL = "https://api.quran.com/api/v3/chapters/%d/info?language=%s"
	wordURL = "https://api.quran.com/api/v3/chapters/%d/verses?language=%s&limit=300&text_type=words"

	languages = []string{
		"en", "ur", "bn", "tr", "es", "ml", "fr", "ru", "bs", "de", "nl", "tg", "id", "it", "ja", "uz", "vi", "zh",
		"sq", "ta", "ms", "bm", "ha", "pt", "ro", "hi", "as", "kk", "sw", "km", "th", "tl", "az", "ko", "ku", "so",
		"bg", "fa", "tt", "zgh", "prs", "am", "ce", "cs", "dv", "fi", "gu", "he", "ka", "kn", "lg", "mk", "mr", "mrn",
		"ne", "no", "om", "pl", "ps", "rw", "sd", "se", "si", "sr", "sq", "sv", "te", "yo", "ug", "uk",
	}
)

func createDownloadRequests() []dl.Request {
	var requests []dl.Request

	for _, lang := range languages {
		// Add surah list URL
		requests = append(requests, dl.Request{
			FileName: fmt.Sprintf("list-%s.json", lang),
			URL:      fmt.Sprintf(listURL, lang),
		})

		// Add surah info URL
		for surah := 1; surah <= 114; surah++ {
			requests = append(requests, dl.Request{
				FileName: fmt.Sprintf("info-%s-%03d.json", lang, surah),
				URL:      fmt.Sprintf(infoURL, surah, lang),
			})
		}
	}

	for _, lang := range languagesForWord {
		// Add word URL
		for surah := 1; surah <= 114; surah++ {
			requests = append(requests, dl.Request{
				FileName: fmt.Sprintf("word-%s-%03d.json", lang, surah),
				URL:      fmt.Sprintf(wordURL, surah, lang),
			})
		}
	}

	return requests
}
