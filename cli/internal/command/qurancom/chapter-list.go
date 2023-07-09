package qurancom

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type LanguagesResponse struct {
	Languages []struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		IsoCode           string `json:"iso_code"`
		NativeName        string `json:"native_name"`
		Direction         string `json:"direction"`
		TranslationsCount int    `json:"translations_count"`
		TranslatedName    struct {
			Name         string `json:"name"`
			LanguageName string `json:"language_name"`
		} `json:"translated_name"`
	} `json:"languages"`
}

type ChaptersResponse struct {
	Chapters []struct {
		ID              int    `json:"id"`
		RevelationPlace string `json:"revelation_place"`
		RevelationOrder int    `json:"revelation_order"`
		BismillahPre    bool   `json:"bismillah_pre"`
		NameSimple      string `json:"name_simple"`
		NameComplex     string `json:"name_complex"`
		NameArabic      string `json:"name_arabic"`
		VersesCount     int    `json:"verses_count"`
		Pages           []int  `json:"pages"`
		TranslatedName  struct {
			LanguageName string `json:"language_name"`
			Name         string `json:"name"`
		} `json:"translated_name"`
	} `json:"chapters"`
}

type ChaptersOutput struct {
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

var (
	urlLanguageList = "https://api.quran.com/api/v4/resources/languages"
	urlChapterList  = "https://api.quran.com/api/v4/chapters?language=%s"
)

func processChapterList(ctx context.Context, cacheDir, dstDir string) error {
	// Download list of available languages
	err := downloadLanguageList(ctx, cacheDir)
	if err != nil {
		return err
	}

	// Parse list of language
	languages, err := parseLanguageList(cacheDir)
	if err != nil {
		return err
	}

	// Download list of chapter names
	err = downloadChapterList(ctx, cacheDir, languages)
	if err != nil {
		return err
	}

	// Save chapter names for each language
	for _, lang := range languages {
		err = parseAndWriteChapterList(cacheDir, dstDir, lang)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadLanguageList(ctx context.Context, cacheDir string) error {
	logrus.Printf("downloading list of languages")
	dstPath := filepath.Join(cacheDir, "list-language.json")

	if !util.FileExist(dstPath) {
		req := dl.Request{URL: urlLanguageList}
		err := dl.Download(ctx, http.DefaultClient, dstPath, req)
		if err != nil {
			return fmt.Errorf("failed to download list language: %w", err)
		}
	}

	return nil
}

func parseLanguageList(cacheDir string) ([]string, error) {
	// Open JSON file
	var listLanguage LanguagesResponse
	listPath := filepath.Join(cacheDir, "list-language.json")
	err := util.DecodeJsonFile(listPath, &listLanguage)
	if err != nil {
		return nil, err
	}

	// We only need the ISO code
	var languages []string
	for _, ch := range listLanguage.Languages {
		languages = append(languages, ch.IsoCode)
	}

	return languages, nil
}

func downloadChapterList(ctx context.Context, cacheDir string, languages []string) error {
	logrus.Printf("downloading list of chapter")

	// Prepare download links
	var dlRequests []dl.Request
	for _, lang := range languages {
		url := fmt.Sprintf(urlChapterList, lang)
		dstName := fmt.Sprintf("chapter-names-%s.json", lang)
		dstPath := filepath.Join(cacheDir, dstName)

		if !util.FileExist(dstPath) {
			dlRequests = append(dlRequests, dl.Request{FileName: dstName, URL: url})
		}
	}

	// Start batch download
	err := dl.BatchDownload(ctx, cacheDir, dlRequests, nil)
	if err != nil {
		return fmt.Errorf("failed to download list of chapter: %w", err)
	}

	return nil
}

func parseAndWriteChapterList(cacheDir, dstDir, lang string) error {
	// Open file
	srcPath := fmt.Sprintf("chapter-names-%s.json", lang)
	srcPath = filepath.Join(cacheDir, srcPath)

	var src ChaptersResponse
	err := util.DecodeJsonFile(srcPath, &src)
	if err != nil {
		return fmt.Errorf("failed to decode chapter list for %s: %w", lang, err)
	}

	// Read each chapter name
	var nMissing int
	result := make(map[string]ChaptersOutput)
	for i, ch := range src.Chapters {
		key := fmt.Sprintf("%03d", i+1)
		name := strings.TrimSpace(ch.NameSimple)
		translation := strings.TrimSpace(ch.TranslatedName.Name)
		if lang != "en" && ch.TranslatedName.LanguageName == "english" {
			translation = "[[MISSING]]"
			nMissing++
		}

		result[key] = ChaptersOutput{
			Name:        name,
			Translation: translation,
		}
	}

	// If many chapter not translated, return empty
	if nMissing > 6 {
		logrus.Warnf("skipped chapter names for %s: %d missing", lang, nMissing)
		return nil
	}

	// Write to file
	logrus.Printf("writing surah names for %s", lang)
	dstDir = filepath.Join(dstDir, "surah-translation")
	os.MkdirAll(dstDir, os.ModePerm)

	dstPath := fmt.Sprintf("%s-qurancom.json", lang)
	dstPath = filepath.Join(dstDir, dstPath)

	// Encode data
	err = util.EncodeSortedKeyJson(dstPath, &result)
	if err != nil {
		return fmt.Errorf("failed to write surah names for %s: %w", lang, err)
	}

	return nil
}
