package islamhouse

import (
	"data-quran-cli/internal/norm"
	"data-quran-cli/internal/util"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-shiori/dom"
	"github.com/sirupsen/logrus"
)

var (
	nAyah = 6_236

	rxNewlines           = regexp.MustCompile(`\s*\n+\s*`)
	rxTafsirNumber       = regexp.MustCompile(`^\(\d+\)\s*`)
	rxTafsirNumberCommon = regexp.MustCompile(`^[\d,\s\.、–-]+\s*`)
	rxTafsirNumberTr     = regexp.MustCompile(`^[\d,\s\.-]+\s+`)
	rxTrHyphenPrefix     = regexp.MustCompile(`^\s*-\s*`)
)

func parseAllPages(cacheDir, language string) ([]string, error) {
	logrus.Printf("parsing page for %s", language)

	// Extract each page in this language
	var tafsirs []string
	for surah := 1; surah <= 114; surah++ {
		sTafsirs, err := parsePage(cacheDir, language, surah)
		if err != nil {
			return nil, err
		}

		tafsirs = append(tafsirs, sTafsirs...)
	}

	// Check if tafsir complete
	if n := len(tafsirs); n != nAyah {
		logrus.Warnf("total tafsir for %s: want %d got %d", language, nAyah, n)
		return nil, nil
	}

	return tafsirs, nil
}

func parsePage(cacheDir, language string, surah int) ([]string, error) {
	// Open file
	path := fmt.Sprintf("%s-mokhtasar-%03d.html", language, surah)
	path = filepath.Join(cacheDir, path)

	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("fail to open page %s %d: %w", language, surah, err)
		return nil, err
	}
	defer f.Close()

	// Parse page to HTML document
	r := norm.NormalizeReader(f)
	doc, err := dom.FastParse(r)
	if err != nil {
		err = fmt.Errorf("fail to parse page %s %d: %w", language, surah, err)
		return nil, err
	}

	// Fetch paragraphs
	var tafsirs []string
	for _, p := range dom.QuerySelectorAll(doc, "#cnt p") {
		text := dom.TextContent(p)
		text = rxNewlines.ReplaceAllString(text, " ")
		text = rxTafsirNumber.ReplaceAllString(text, "")

		// Do some little cleaning
		switch language {
		case "ja": // Do nothing
		case "tr":
			text = rxTrHyphenPrefix.ReplaceAllString(text, "-")
			text = rxTafsirNumberTr.ReplaceAllString(text, "")
		default:
			text = rxTafsirNumberCommon.ReplaceAllString(text, "")
		}

		text = util.MarkdownText(text)
		tafsirs = append(tafsirs, text)
	}

	// Check count
	nTafsir := len(tafsirs)
	nExpected := util.ListSurah[surah].NAyah
	if nTafsir != nExpected {
		err = fmt.Errorf("page %s %d: want %d got %d ayah", language, surah, nExpected, nTafsir)
		return nil, err
	}

	return tafsirs, nil
}
