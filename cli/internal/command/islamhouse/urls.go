package islamhouse

import (
	"data-quran-cli/internal/dl"
	"fmt"
)

type SourceName struct {
	Language string
	Slug     string
}

var (
	baseURL = "https://islamhouse.com/quran/%s/sura-%d.html"

	sourceNames []SourceName = []SourceName{
		{"ar", "arabic_mokhtasar"},
		{"tr", "turkish_mokhtasar"},
		{"fr", "french_mokhtasar"},
		{"id", "indonesian_mokhtasar"},
		{"vi", "vietnamese_mokhtasar"},
		{"bs", "bosnian_mokhtasar"},
		{"it", "italian_mokhtasar"},
		{"es", "spanish_mokhtasar"},
		{"tl", "tagalog_mokhtasar"},
		{"bn", "bengali_mokhtasar"},
		{"fa", "persian_mokhtasar"},
		{"zh", "chinese_mokhtasar"},
		{"ja", "japanese_mokhtasar"},
		{"as", "assamese_mokhtasar"},
		{"ml", "malayalam_mokhtasar"},
		{"km", "khmer_mokhtasar"},
	}
)

func createDownloadRequests() []dl.Request {
	var requests []dl.Request

	for _, path := range sourceNames {
		for surah := 1; surah <= 114; surah++ {
			requests = append(requests, dl.Request{
				FileName: fmt.Sprintf("%s-mokhtasar-%03d.html", path.Language, surah),
				URL:      fmt.Sprintf(baseURL, path.Slug, surah),
			})
		}
	}

	return requests
}
