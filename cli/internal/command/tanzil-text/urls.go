package tanzilText

import (
	"data-quran-cli/internal/dl"
	"net/url"
)

var quranTypes = []string{
	"simple",
	"simple-plain",
	"simple-min",
	"simple-clean",
	"uthmani",
	"uthmani-min",
}

func createQuranURLs() []dl.Request {
	// Prepare base URL
	query := url.Values{
		"marks":   {"true"},
		"sajdah":  {"true"},
		"rub":     {"true"},
		"tatweel": {"true"},
		"outType": {"txt"},
		"agree":   {"true"},
	}

	baseURL := url.URL{
		Scheme: "https",
		Host:   "tanzil.net",
		Path:   "pub/download/index.php",
	}

	// Generate download requests
	var requests []dl.Request
	for _, tp := range quranTypes {
		query.Set("quranType", tp)
		baseURL.RawQuery = query.Encode()
		requests = append(requests, dl.Request{
			URL:      baseURL.String(),
			FileName: tp + ".txt",
		})
	}

	return requests
}
