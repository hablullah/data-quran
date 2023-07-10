package qurancom

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	iso6391 "github.com/emvi/iso-639-1"
	"github.com/go-shiori/dom"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
)

type ListTafsirResponseEntry struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	AuthorName     string `json:"author_name"`
	Slug           string `json:"slug"`
	LanguageName   string `json:"language_name"`
	TranslatedName struct {
		Name         string `json:"name"`
		LanguageName string `json:"language_name"`
	} `json:"translated_name"`
}

type ListTafsirResponse struct {
	Tafsirs []ListTafsirResponseEntry `json:"tafsirs"`
}

type TafsirResponse struct {
	Tafsirs []struct {
		ID         int    `json:"id"`
		ResourceID int    `json:"resource_id"`
		VerseKey   string `json:"verse_key"`
		LanguageID int    `json:"language_id"`
		Text       string `json:"text"`
		Slug       string `json:"slug"`
	} `json:"tafsirs"`
	Pagination struct {
		PerPage      int `json:"per_page"`
		CurrentPage  int `json:"current_page"`
		NextPage     any `json:"next_page"`
		TotalPages   int `json:"total_pages"`
		TotalRecords int `json:"total_records"`
	} `json:"pagination"`
}

var (
	urlTafsirList  = "https://api.quran.com/api/v4/resources/tafsirs"
	urlTafsirSurah = "https://api.quran.com/api/v4/tafsirs/%d/by_chapter/%d?per_page=300"

	rxRoundBrackets = regexp.MustCompile(`\s*\([^)]*\)`)
)

func downloadTafsirList(ctx context.Context, cacheDir string) error {
	logrus.Printf("downloading list of tafsirs")
	dstPath := filepath.Join(cacheDir, "list-tafsir.json")

	if !util.FileExist(dstPath) {
		req := dl.Request{URL: urlTafsirList}
		err := dl.Download(ctx, http.DefaultClient, dstPath, req)
		if err != nil {
			return fmt.Errorf("failed to download list of tafsir: %w", err)
		}
	}

	return nil
}

func parseTafsirList(cacheDir string) ([]ListTafsirResponseEntry, error) {
	// Open JSON file
	var listTafsir ListTafsirResponse
	listPath := filepath.Join(cacheDir, "list-tafsir.json")
	err := util.DecodeJsonFile(listPath, &listTafsir)
	if err != nil {
		return nil, err
	}

	return listTafsir.Tafsirs, nil
}

func downloadTafsir(ctx context.Context, cacheDir string, entry ListTafsirResponseEntry) error {
	logrus.Printf("downloading tafsir %q", entry.Slug)

	// Prepare download links
	var dlRequests []dl.Request
	for surah := 1; surah <= 114; surah++ {
		url := fmt.Sprintf(urlTafsirSurah, entry.ID, surah)
		dstName := fmt.Sprintf("tafsir-%s-%03d.json", entry.Slug, surah)
		dstPath := filepath.Join(cacheDir, dstName)

		if !util.FileExist(dstPath) {
			dlRequests = append(dlRequests, dl.Request{FileName: dstName, URL: url})
		}
	}

	// Start batch download
	err := dl.BatchDownload(ctx, cacheDir, dlRequests, nil)
	if err != nil {
		return fmt.Errorf("failed to download tafsir %q: %w", entry.Slug, err)
	}

	return nil
}

func downloadAllTafsirs(ctx context.Context, cacheDir string, entries []ListTafsirResponseEntry) error {
	for _, t := range entries {
		err := downloadTafsir(ctx, cacheDir, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseTafsir(cacheDir string, entry ListTafsirResponseEntry) ([]string, error) {
	logrus.Printf("parsing tafsir %q", entry.Slug)

	// Extract tafsirs
	var tafsirs []string
	for surah := 1; surah <= 114; surah++ {
		// Decode JSON
		var data TafsirResponse
		srcPath := fmt.Sprintf("tafsir-%s-%03d.json", entry.Slug, surah)
		srcPath = filepath.Join(cacheDir, srcPath)
		err := util.DecodeJsonFile(srcPath, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to decode tafsir %q:%d: %w", entry.Slug, surah, err)
		}

		// Process each ayah
		for _, t := range data.Tafsirs {
			// Parse document
			doc, err := dom.FastParse(strings.NewReader(t.Text))
			if err != nil {
				return nil, fmt.Errorf("failed to parse tafsir %q:%d: %w", entry.Slug, surah, err)
			}

			// Replace all H1 to H2, and so on
			if len(dom.GetElementsByTagName(doc, "h1")) > 0 {
				for hLevel := 5; hLevel >= 1; hLevel-- {
					tagName := fmt.Sprintf("h%d", hLevel)
					newTagName := fmt.Sprintf("h%d", hLevel+1)
					hNodes := dom.GetElementsByTagName(doc, tagName)
					for _, node := range hNodes {
						node.Data = newTagName
					}
				}
			}

			// Return back doc to text
			body := dom.QuerySelector(doc, "body")
			text := dom.InnerHTML(body)
			text = util.MarkdownText(text)

			// Save text to slice
			tafsirs = append(tafsirs, text)
		}
	}

	return tafsirs, nil
}

func writeTafsir(dstDir string, listAyah []string, entry ListTafsirResponseEntry) error {
	logrus.Printf("writing tafsir for %q", entry.Slug)

	// Make sure count of ayah is correct
	if nAyah := len(listAyah); nAyah != 6236 {
		logrus.Warnf("%q n ayah %d != 6236", entry.Slug, nAyah)
		return nil
	}

	// Prepare metadata
	name := rxRoundBrackets.ReplaceAllString(entry.TranslatedName.Name, "")
	name = strings.TrimSpace(name)
	author := strings.TrimSpace(entry.AuthorName)
	language := util.UpperFirst(strings.TrimSpace(entry.LanguageName))
	slug := strings.TrimSpace(entry.Slug)
	id := entry.ID

	// Prepare filename
	isoLang := iso6391.CodeForName(language)
	if isoLang == "" {
		return fmt.Errorf("unknown language %q in tafsir %q", language, slug)
	}

	dstName := strings.ToLower(name)
	dstName = strings.ReplaceAll(dstName, "tafsir", "")
	dstName = strings.ReplaceAll(dstName, "tafseer", "")
	dstName = strings.TrimSpace(dstName)
	dstName = strcase.ToKebab(dstName)
	dstName = fmt.Sprintf("%s-%s-qurancom.md", isoLang, dstName)
	dstPath := filepath.Join(dstDir, "ayah-tafsir", dstName)

	// Write metadata
	var sb strings.Builder
	sb.WriteString("<!--\n")
	sb.WriteString(fmt.Sprintf("Name    : %s\n", name))
	sb.WriteString(fmt.Sprintf("Author  : %s\n", author))
	sb.WriteString(fmt.Sprintf("Language: %s\n", language))
	sb.WriteString(fmt.Sprintf("Source  : %s\n", "Quran.com"))
	sb.WriteString(fmt.Sprintf("Slug    : %s\n", slug))
	sb.WriteString(fmt.Sprintf("ID      : %d\n", id))
	sb.WriteString("-->\n\n")

	// Write each ayah
	for i, ayah := range listAyah {
		sb.WriteString("# ")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString("\n\n")

		if ayah == "" {
			sb.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		}

		if i > 0 {
			prevAyah := listAyah[i-1]
			if ayah == prevAyah {
				sb.WriteString("<!-- TODO:DUPLICATE -->\n\n")
			}
		}

		sb.WriteString(ayah)
		sb.WriteString("\n\n")
	}

	// Write to file
	os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	err := os.WriteFile(dstPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write tafsir %q: %w", slug, err)
	}

	return nil
}
