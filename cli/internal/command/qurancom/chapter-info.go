package qurancom

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	iso6391 "github.com/emvi/iso-639-1"
	"github.com/go-shiori/dom"
	"github.com/sirupsen/logrus"
)

type ChapterInfosResponse struct {
	ChapterInfos []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		AuthorName     string `json:"author_name"`
		Slug           string `json:"slug"`
		LanguageName   string `json:"language_name"`
		TranslatedName struct {
			Name         string `json:"name"`
			LanguageName string `json:"language_name"`
		} `json:"translated_name"`
	} `json:"chapter_infos"`
}

type ChapterInfoResponse struct {
	ChapterInfo struct {
		ID           int    `json:"id"`
		ChapterID    int    `json:"chapter_id"`
		LanguageName string `json:"language_name"`
		ShortText    string `json:"short_text"`
		Source       string `json:"source"`
		Text         string `json:"text"`
	} `json:"chapter_info"`
}

type ChapterInfoOutput struct {
	Number   int
	Language string
	Source   string
	Text     string
}

var (
	urlChapterInfoList = "https://api.quran.com/api/v4/resources/chapter_infos"
	urlChapterInfo     = "https://api.quran.com/api/v4/chapters/%d/info?language=%s"
)

func processChapterInfo(ctx context.Context, cacheDir, dstDir string) error {
	// Download list of chapter info
	err := downloadChapterInfoList(ctx, cacheDir)
	if err != nil {
		return err
	}

	// Parse list of chapter info
	languages, err := parseChapterInfoList(cacheDir)
	if err != nil {
		return err
	}

	// Download chapter info for each language
	err = downloadAllChapterInfo(ctx, cacheDir, languages)
	if err != nil {
		return err
	}

	// Parse and write each chapter info
	for _, lang := range languages {
		results, err := parseAllChapterInfo(cacheDir, lang)
		if err != nil {
			return err
		}

		err = writeChapterInfo(dstDir, lang, results)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadChapterInfoList(ctx context.Context, cacheDir string) error {
	logrus.Printf("downloading list of chapter info")
	dstPath := filepath.Join(cacheDir, "list-chapter-info.json")

	if !util.FileExist(dstPath) {
		req := dl.Request{URL: urlChapterInfoList}
		err := dl.Download(ctx, http.DefaultClient, dstPath, req)
		if err != nil {
			return fmt.Errorf("failed to download list chapter info: %w", err)
		}
	}

	return nil
}

func parseChapterInfoList(cacheDir string) ([]string, error) {
	// Open JSON file
	var listChapterInfo ChapterInfosResponse
	listPath := filepath.Join(cacheDir, "list-chapter-info.json")
	err := util.DecodeJsonFile(listPath, &listChapterInfo)
	if err != nil {
		return nil, err
	}

	// Convert languages to ISO-639-1 code
	var languages []string
	for _, ch := range listChapterInfo.ChapterInfos {
		lang := util.UpperFirst(ch.LanguageName)
		isoLang := iso6391.CodeForName(lang)
		if isoLang == "" {
			return nil, fmt.Errorf("unknown language: %q", lang)
		}

		languages = append(languages, isoLang)
	}

	return languages, nil
}

func downloadChapterInfo(ctx context.Context, cacheDir string, isoLang string) error {
	logrus.Printf("downloading chapter info for %s", isoLang)

	// Prepare download links
	var dlRequests []dl.Request
	for surah := 1; surah <= 114; surah++ {
		url := fmt.Sprintf(urlChapterInfo, surah, isoLang)
		dstName := fmt.Sprintf("chapter-info-%s-%03d.json", isoLang, surah)
		dstPath := filepath.Join(cacheDir, dstName)

		if !util.FileExist(dstPath) {
			dlRequests = append(dlRequests, dl.Request{FileName: dstName, URL: url})
		}
	}

	// Start batch download
	err := dl.BatchDownload(ctx, cacheDir, dlRequests, nil)
	if err != nil {
		return fmt.Errorf("failed to download %s chapter info: %w", isoLang, err)
	}

	return nil
}

func downloadAllChapterInfo(ctx context.Context, cacheDir string, languages []string) error {
	for _, lang := range languages {
		err := downloadChapterInfo(ctx, cacheDir, lang)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseChapterInfo(cacheDir, lang string, surah int) (*ChapterInfoOutput, error) {
	// Open file
	srcPath := fmt.Sprintf("chapter-info-%s-%03d.json", lang, surah)
	srcPath = filepath.Join(cacheDir, srcPath)

	var src ChapterInfoResponse
	err := util.DecodeJsonFile(srcPath, &src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode chapter info for %s %d: %w", lang, surah, err)
	}

	// If language is not english, but we got one in neglish, remove it
	srcText := src.ChapterInfo.Text
	if lang != "en" && src.ChapterInfo.LanguageName == "english" {
		srcText = ""
	}

	// If text is empty, just stop
	if srcText == "" {
		return nil, nil
	}

	// Convert text to html.Node document
	doc, err := dom.FastParse(strings.NewReader(srcText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse chapter info for %s %d: %w", lang, surah, err)
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
	docHTML := dom.InnerHTML(doc)

	// Return output
	return &ChapterInfoOutput{
		Text:     util.MarkdownText(docHTML),
		Number:   src.ChapterInfo.ChapterID,
		Source:   src.ChapterInfo.Source,
		Language: src.ChapterInfo.LanguageName,
	}, nil
}

func parseAllChapterInfo(cacheDir, lang string) ([]*ChapterInfoOutput, error) {
	var outputs []*ChapterInfoOutput
	for surah := 1; surah <= 114; surah++ {
		output, err := parseChapterInfo(cacheDir, lang, surah)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}

func writeChapterInfo(dstDir string, language string, data []*ChapterInfoOutput) error {
	// If data is empty, stop
	if len(data) != 114 {
		return nil
	}

	logrus.Printf("writing surah info for %q", language)

	// Write metadata
	var sb strings.Builder
	sb.WriteString("<!--\n")
	sb.WriteString(fmt.Sprintf("Language: %s\n", data[0].Language))
	sb.WriteString(fmt.Sprintf("Source  : %s\n", data[0].Source))
	sb.WriteString("-->\n\n")

	// Write each info
	for i, info := range data {
		sb.WriteString("# ")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString("\n\n")

		if info == nil || info.Text == "" {
			sb.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		}

		sb.WriteString(info.Text)
		sb.WriteString("\n\n")
	}

	// Prepare destination path
	dstDir = filepath.Join(dstDir, "surah-info")
	os.MkdirAll(dstDir, os.ModePerm)

	dstPath := fmt.Sprintf("%s-qurancom.md", language)
	dstPath = filepath.Join(dstDir, dstPath)

	// Write the file
	err := os.WriteFile(dstPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write %s surah info: %w", language, err)
	}

	return nil
}
