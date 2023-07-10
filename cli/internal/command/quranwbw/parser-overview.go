package quranwbw

import (
	"data-quran-cli/internal/util"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var rxOverviewH2 = regexp.MustCompile(`(?i)<h2>[^<]+</h2>`)

func parseOverviews(cacheDir string) (map[int]string, error) {
	// Open file
	var overviews map[int]string
	srcPath := filepath.Join(cacheDir, "000-overview.json")
	err := util.DecodeJsonFile(srcPath, &overviews)
	if err != nil {
		return nil, fmt.Errorf("failed to decode overview: %w", err)
	}

	// Process and normalize data
	for surah, overview := range overviews {
		overview = rxOverviewH2.ReplaceAllString(overview, "")
		overview = strings.ReplaceAll(overview, "h3", "h2")
		overview = util.MarkdownText(overview)
		overviews[surah] = overview
	}

	return overviews, nil
}
