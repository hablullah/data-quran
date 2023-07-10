package kemenag

import (
	"data-quran-cli/internal/util"
	"encoding/json"
	"fmt"
	"net/http"
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
	urlQuranTafsir = "https://web-api.qurankemenag.net/quran-tafsir/%d"

	rxNewlines        = regexp.MustCompile(`\n+`)
	rxTafsirNumber    = regexp.MustCompile(`(?m)^\s*([\d\-]+)\s*\\?\.?\s*`)
	rxTrimTafsirSpace = regexp.MustCompile(`(?m)^[^\S\n]+|[^\S\n]+$`)
	rxWajizBracket    = regexp.MustCompile(`\s*\\\[\s*\\\]\s*`)

	rxTransAyahNumber = regexp.MustCompile(`^\s*([\d\-]+)\s*\\?\.?\s*`)
	rxTransFnNumber   = regexp.MustCompile(`\s*(\d+)\s*\)(\s*)`)
	rxFootFnNumber    = regexp.MustCompile(`(?m)^\s*(\d+)\s*\)\s*`)
)

func processTafsir(cacheDir, dstDir string) error {
	// Download tafsirs
	err := downladAllTafsir(cacheDir)
	if err != nil {
		return err
	}

	// Parse each surah to extract text, trans and tafsirs
	listAyah, err := parseAllTafsir(cacheDir)
	if err != nil {
		return err
	}

	// Save surah data
	err = writeQuranBasicData(dstDir, listAyah, TextArabic)
	if err != nil {
		return err
	}

	err = writeQuranBasicData(dstDir, listAyah, Transliteration)
	if err != nil {
		return err
	}

	err = writeQuranBasicData(dstDir, listAyah, TafsirWajiz)
	if err != nil {
		return err
	}

	err = writeQuranBasicData(dstDir, listAyah, TafsirTahlili)
	if err != nil {
		return err
	}

	err = writeQuranTranslation(dstDir, listAyah)
	if err != nil {
		return err
	}

	err = writeSurahInfo(dstDir, listAyah)
	if err != nil {
		return err
	}

	return nil
}

func downladAllTafsir(cacheDir string) error {
	for surah := 1; surah <= 114; surah++ {
		err := downloadTafsir(cacheDir, surah)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadTafsir(cacheDir string, surah int) error {
	// Prepare destination path
	dstName := fmt.Sprintf("surah-%03d.json", surah)
	dstPath := filepath.Join(cacheDir, dstName)
	if util.FileExist(dstPath) {
		return nil
	}

	// Prepare http client
	client := &http.Client{}

	// Download each tafsir
	surahData := util.ListSurah[surah]
	tafsirs := make([]Ayah, surahData.NAyah)

	for idx := 1; idx <= surahData.NAyah; idx++ {
		ayah := surahData.Start + idx - 1
		err := func() error {
			logrus.Printf("downloading tafsir for %d:%d", surah, idx)

			// Download page
			url := fmt.Sprintf(urlQuranTafsir, ayah)
			resp, err := client.Get(url)
			if err != nil {
				return fmt.Errorf("failed to download tafsir for %d:%d, %w", surah, idx, err)
			}
			defer resp.Body.Close()

			// Decode data
			var respData RespDownloadTafsir
			err = json.NewDecoder(resp.Body).Decode(&respData)
			if err != nil {
				return fmt.Errorf("failed to decode tafsir for %d:%d %w", surah, idx, err)
			}

			// Save to slice
			tafsirs[idx-1] = respData.Data
			return nil
		}()
		if err != nil {
			return err
		}
	}

	// Write tafsirs to file
	bt, err := json.MarshalIndent(tafsirs, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to save tafsir for surah %d: %w", surah, err)
	}

	return os.WriteFile(dstPath, bt, os.ModePerm)
}

func parseAllTafsir(cacheDir string) ([]Ayah, error) {
	allAyah := make([]Ayah, 0, 6236)
	for surah := 1; surah <= 114; surah++ {
		listAyah, err := parseTafsirInSurah(cacheDir, surah)
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

func parseTafsirInSurah(cacheDir string, surah int) ([]Ayah, error) {
	logrus.Printf("parsing surah %d", surah)

	// Prepare path
	srcPath := fmt.Sprintf("surah-%03d.json", surah)
	srcPath = filepath.Join(cacheDir, srcPath)

	// Decode JSON
	var listAyah []Ayah
	err := util.DecodeJsonFile(srcPath, &listAyah)
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
	os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
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

func writeSurahInfo(dstDir string, listAyah []Ayah) error {
	logrus.Printf("writing surah info")

	fnClean := func(s string) string {
		s = util.MarkdownText(s)
		s = rxTafsirNumber.ReplaceAllString(s, "${1}. ")
		s = rxTrimTafsirSpace.ReplaceAllString(s, "")
		s = rxNewlines.ReplaceAllString(s, "\n\n")
		return s
	}

	// Write metadata
	var sb strings.Builder
	sb.WriteString("<!--\n")
	sb.WriteString("Title : Penjelasan Surah dari Quran Kemenag\n")
	sb.WriteString("Source: https://quran.kemenag.go.id\n")
	sb.WriteString("-->\n\n")

	// Write each ayah
	for surah := 1; surah <= 114; surah++ {
		// Fetch data
		surahData := util.ListSurah[surah]
		lastAyahIdx := surahData.Start + surahData.NAyah - 1 - 1
		lastAyah := listAyah[lastAyahIdx]
		outro := fnClean(lastAyah.Tafsir.OutroSurah.String)
		intro := fnClean(lastAyah.Tafsir.IntroSurah.String)

		var surahRelation string
		firstAyahOfNextSurahIdx := lastAyahIdx + 1
		if firstAyahOfNextSurahIdx < 6236 {
			firstAyahOfNextSurah := listAyah[firstAyahOfNextSurahIdx]
			surahRelation = firstAyahOfNextSurah.Tafsir.MunasabahPrevSurah.String
			surahRelation = fnClean(surahRelation)
		} else {
			surahRelation = lastAyah.Tafsir.MunasabahPrevSurah.String
			surahRelation = fnClean(surahRelation)
		}

		// Merge info into one
		var info []string
		if intro != "" {
			info = append(info, intro)
		}
		if outro != "" {
			info = append(info, outro)
		}
		if surahRelation != "" {
			info = append(info, surahRelation)
		}

		// Write info
		sb.WriteString("# ")
		sb.WriteString(strconv.Itoa(surah))
		sb.WriteString("\n\n")

		if len(info) == 0 {
			sb.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		} else {
			sb.WriteString(strings.Join(info, "\n\n"))
			sb.WriteString("\n\n")
		}
	}

	// Write to file
	dstPath := filepath.Join(dstDir, "surah-info", "id-kemenag.md")
	err := os.WriteFile(dstPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write info: %w", err)
	}

	return nil
}
