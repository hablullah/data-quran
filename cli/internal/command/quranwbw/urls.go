package quranwbw

import (
	"data-quran-cli/internal/dl"
	"fmt"
)

var languages = map[string]string{
	"arabic":     "ar",
	"english":    "en",
	"urdu":       "ur",
	"hindi":      "hi",
	"indonesian": "id",
	"bangla":     "bn",
	"turkish":    "tr",
	"german":     "de",
	"russian":    "ru",
	"ingush":     "inh",
	"tamil":      "ta",
}

func createDownloadRequests() []dl.Request {
	var requests []dl.Request
	baseURL := "https://data.quranwbw.com/%d/word-translations/%s.json?v1680114261"

	for lang := range languages {
		for surah := 1; surah <= 114; surah++ {
			requests = append(requests, dl.Request{
				URL:      fmt.Sprintf(baseURL, surah, lang),
				FileName: fmt.Sprintf("%s-%03d.json", lang, surah),
			})
		}
	}

	return requests
}
