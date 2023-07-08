package kemenag

import (
	"data-quran-cli/internal/norm"
	"data-quran-cli/internal/util"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type BasicDataField uint8

func (f BasicDataField) String() string {
	switch f {
	case TextArabic:
		return "Text Arabic"
	case Transliteration:
		return "Transliteration"
	case TafsirWajiz:
		return "Tafsir Wajiz"
	case TafsirTahlili:
		return "Tafsir Tahlili"
	default:
		return ""
	}
}

const (
	TextArabic BasicDataField = iota + 1
	Transliteration
	TafsirWajiz
	TafsirTahlili
)

var (
	rxNewlines        = regexp.MustCompile(`\n+`)
	rxTafsirNumber    = regexp.MustCompile(`(?m)^\s*([\d\-]+)\s*\\?\.?\s*`)
	rxTrimTafsirSpace = regexp.MustCompile(`(?m)^[^\S\n]+|[^\S\n]+$`)
	rxWajizBracket    = regexp.MustCompile(`\s*\\\[\s*\\\]\s*`)

	rxTransAyahNumber = regexp.MustCompile(`^\s*([\d\-]+)\s*\\?\.?\s*`)
	rxTransFnNumber   = regexp.MustCompile(`\s*(\d+)\s*\)(\s*)`)
	rxFootFnNumber    = regexp.MustCompile(`(?m)^\s*(\d+)\s*\)\s*`)
)

func parseAllSurah(cacheDir string) ([]Ayah, error) {
	allAyah := make([]Ayah, 0, 6236)
	for surah := 1; surah <= 114; surah++ {
		listAyah, err := parseSurah(cacheDir, surah)
		if err != nil {
			return nil, err
		}
		allAyah = append(allAyah, listAyah...)
	}

	// Make sure there are 6236 ayah
	if nAyah := len(allAyah); nAyah != 6236 {
		return nil, fmt.Errorf("n ayah %d != 6236", nAyah)
	}

	return allAyah, nil
}

func parseSurah(cacheDir string, surah int) ([]Ayah, error) {
	logrus.Printf("parsing surah %d", surah)

	// Open file
	srcPath := fmt.Sprintf("surah-%03d.json", surah)
	srcPath = filepath.Join(cacheDir, srcPath)
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read surah %d: %w", surah, err)
	}
	defer f.Close()

	// Decode JSON
	var listAyah []Ayah
	r := norm.NormalizeReader(f)
	err = json.NewDecoder(r).Decode(&listAyah)
	if err != nil {
		return nil, fmt.Errorf("failed to decode surah %d: %w", surah, err)
	}

	return listAyah, nil
}

func writeQuranBasicData(dstDir string, listAyah []Ayah, field BasicDataField) error {
	logrus.Printf("writing quran data: %s", field)

	// Prepare data
	var title string
	var dstPath string
	var fnValue func(Ayah) string
	fnCleanTafsir := func(s string) string {
		s = util.MarkdownText(s)
		s = rxTafsirNumber.ReplaceAllString(s, "${1}. ")
		s = rxTrimTafsirSpace.ReplaceAllString(s, "")
		s = rxNewlines.ReplaceAllString(s, "\n\n")
		return s
	}

	switch field {
	case TextArabic:
		title = "Quran Kemenag"
		fnValue = func(a Ayah) string { return util.MarkdownText(a.Arabic) }
		dstPath = filepath.Join(dstDir, "ayah-text", "mishbah-kemenag.md")
	case Transliteration:
		title = "Quran Kemenag"
		fnValue = func(a Ayah) string { return util.MarkdownText(a.Latin) }
		dstPath = filepath.Join(dstDir, "ayah-transliteration", "id-transliteration-kemenag.md")
	case TafsirWajiz:
		title = "Tafsir Wajiz (Ringkas) Kemenag"
		dstPath = filepath.Join(dstDir, "ayah-tafsir", "id-wajiz-kemenag.md")
		fnValue = func(a Ayah) string {
			s := fnCleanTafsir(a.Tafsir.Wajiz)
			s = rxWajizBracket.ReplaceAllString(s, "")
			return s
		}
	case TafsirTahlili:
		title = "Tafsir Tahlili Kemenag"
		dstPath = filepath.Join(dstDir, "ayah-tafsir", "id-tahlili-kemenag.md")
		fnValue = func(a Ayah) string {
			s := fnCleanTafsir(a.Tafsir.Tahlili)
			return s
		}
	}

	// Write metadata
	var sb strings.Builder
	sb.WriteString("<!--\n")
	sb.WriteString("Title : " + title + "\n")
	sb.WriteString("Source: https://quran.kemenag.go.id\n")
	if field == TextArabic {
		sb.WriteString("Best used with font LPMQ Isep Mishbah: https://lajnah.kemenag.go.id/unduhan/quran-kemenag.html\n")
	}
	sb.WriteString("-->\n\n")

	// Write each ayah
	for i, ayah := range listAyah {
		sb.WriteString("# ")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString("\n\n")

		value := strings.TrimSpace(fnValue(ayah))
		if value == "" {
			sb.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		}

		if i > 0 {
			prevAyah := listAyah[i-1]
			prevValue := strings.TrimSpace(fnValue(prevAyah))
			if value == prevValue {
				sb.WriteString("<!-- TODO:DUPLICATE -->\n\n")
			}
		}

		sb.WriteString(value)
		sb.WriteString("\n\n")
	}

	// Write to file
	err := os.WriteFile(dstPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write data %q: %w", field, err)
	}

	return nil
}

func writeQuranTranslation(dstDir string, listAyah []Ayah) error {
	logrus.Printf("writing quran translation")

	// Write metadata
	var sb strings.Builder
	sb.WriteString("<!--\n")
	sb.WriteString("Title : Terjemah Quran Kemenag\n")
	sb.WriteString("Source: https://quran.kemenag.go.id\n")
	sb.WriteString("-->\n\n")

	// Write each ayah
	for i, ayah := range listAyah {
		sb.WriteString("# ")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString("\n\n")

		// Check translation
		translation := strings.TrimSpace(ayah.Translation)
		translation = util.MarkdownText(translation)
		translation = rxTransAyahNumber.ReplaceAllString(translation, "")
		translation = rxNewlines.ReplaceAllString(translation, " ")

		if translation == "" {
			sb.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		}

		if i > 0 {
			prevAyah := listAyah[i-1]
			prevTranslation := strings.TrimSpace(prevAyah.Translation)
			if translation == prevTranslation {
				sb.WriteString("<!-- TODO:DUPLICATE -->\n\n")
			}
		}

		// Prepare footnote
		footnote := strings.TrimSpace(ayah.Footnotes.String)
		footnote = util.MarkdownText(footnote)

		// If there are no footnote, continue to the next
		if footnote == "" {
			sb.WriteString(translation)
			sb.WriteString("\n\n")
		} else {
			translation = rxTransFnNumber.ReplaceAllString(translation, "[^${1}]${2}")
			translation = strings.ReplaceAll(translation, "\\[[^", "[^")
			footnote = rxFootFnNumber.ReplaceAllString(footnote, "\n\n[^${1}]: ")
			footnote = rxNewlines.ReplaceAllString(footnote, "\n\n")
			footnote = strings.TrimSpace(footnote)

			sb.WriteString(translation)
			sb.WriteString("\n\n")
			sb.WriteString(footnote)
			sb.WriteString("\n\n")
		}
	}

	// Write to file
	dstPath := filepath.Join(dstDir, "ayah-translation", "id-kemenag.md")
	err := os.WriteFile(dstPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write translation: %w", err)
	}

	return nil
}
