package quranwbw

import (
	"data-quran-cli/internal/norm"
	"data-quran-cli/internal/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rxOverviewH2 = regexp.MustCompile(`(?i)<h2>[^<]+</h2>`)

func parseOverviews(cacheDir string) (map[int]string, error) {
	// Open file
	srcPath := filepath.Join(cacheDir, "000-overview.json")
	src, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open overview: %w", err)
	}
	defer src.Close()

	// Decode source file
	var overviews map[int]string
	err = json.NewDecoder(src).Decode(&overviews)
	if err != nil {
		return nil, fmt.Errorf("failed to decode overview: %w", err)
	}

	// Process and normalize data
	for surah, overview := range overviews {
		overview = norm.NormalizeUnicode(overview)
		overview = rxOverviewH2.ReplaceAllString(overview, "")
		overview = strings.ReplaceAll(overview, "h3", "h2")
		overview = util.MarkdownText(overview)
		overviews[surah] = overview
	}

	return overviews, nil
}
