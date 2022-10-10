package kemenag

import (
	"data-quran-cli/internal/dl"
	"fmt"
)

func createDownloadRequests() []dl.Request {
	var requests []dl.Request

	// Add list surah
	requests = append(requests, dl.Request{
		FileName: "list-surah.json",
		URL:      "https://quran.kemenag.go.id/api/v1/surah/0/114",
	})

	// Add translation URLs
	transURL := "https://quran.kemenag.go.id/api/v1/ayatweb/%d/0/0/300"
	for surah := 1; surah <= 114; surah++ {
		requests = append(requests, dl.Request{
			FileName: fmt.Sprintf("surah-%03d.json", surah),
			URL:      fmt.Sprintf(transURL, surah),
		})
	}

	// Add tafsir URLs
	tafsirURL := "https://quran.kemenag.go.id/api/v1/tafsirbyayat/%d"
	for ayah := 1; ayah <= 6236; ayah++ {
		requests = append(requests, dl.Request{
			FileName: fmt.Sprintf("ayah-%04d.json", ayah),
			URL:      fmt.Sprintf(tafsirURL, ayah),
		})
	}

	return requests
}
