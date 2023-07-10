package qurancom

import (
	"context"
	"data-quran-cli/internal/dl"
	"data-quran-cli/internal/util"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
)

type AyahTextResponse struct {
	Verses []AyahTextResponseVerse `json:"verses"`
	Meta   struct {
		Filters struct {
			ChapterNumber string `json:"chapter_number"`
		} `json:"filters"`
	} `json:"meta"`
}

type AyahTextResponseVerse struct {
	ID                 int    `json:"id"`
	VerseKey           string `json:"verse_key"`
	TextImlaei         string `json:"text_imlaei"`
	TextImlaeiSimple   string `json:"text_imlaei_simple"`
	TextIndopak        string `json:"text_indopak"`
	TextUthmani        string `json:"text_uthmani"`
	TextUthmaniSimple  string `json:"text_uthmani_simple"`
	TextUthmaniTajweed string `json:"text_uthmani_tajweed"`
}

var (
	ayahTextNames = []string{
		"imlaei_simple",
		"imlaei",
		"indopak",
		"uthmani_simple",
		"uthmani_tajweed",
		"uthmani",
	}
)

func downloadAllAyahText(ctx context.Context, cacheDir string) error {
	logrus.Printf("downloading all ayah text")

	// Prepare download links
	var dlRequests []dl.Request
	for _, name := range ayahTextNames {
		url := fmt.Sprintf("https://api.quran.com/api/v4/quran/verses/%s", name)
		dstName := strcase.ToKebab(name)
		dstName = fmt.Sprintf("ayah-text-%s.json", dstName)
		dstPath := filepath.Join(cacheDir, dstName)

		if !util.FileExist(dstPath) {
			dlRequests = append(dlRequests, dl.Request{FileName: dstName, URL: url})
		}
	}

	// Start batch download
	err := dl.BatchDownload(ctx, cacheDir, dlRequests, nil)
	if err != nil {
		return fmt.Errorf("failed to download ayah text: %w", err)
	}

	return nil
}

func parseAndWriteAyahText(cacheDir, dstDir, name string) error {
	logrus.Printf("writing ayah text: %s", name)

	// Open JSON file
	var ayahText AyahTextResponse
	srcPath := strcase.ToKebab(name)
	srcPath = fmt.Sprintf("ayah-text-%s.json", srcPath)
	srcPath = filepath.Join(cacheDir, srcPath)
	err := util.DecodeJsonFile(srcPath, &ayahText)
	if err != nil {
		return fmt.Errorf("failed to decode %s text: %w", name, err)
	}

	// Prepare helper function
	var fnValue func(v AyahTextResponseVerse) string
	switch name {
	case "imlaei":
		fnValue = func(v AyahTextResponseVerse) string { return v.TextImlaei }
	case "imlaei_simple":
		fnValue = func(v AyahTextResponseVerse) string { return v.TextImlaeiSimple }
	case "indopak":
		fnValue = func(v AyahTextResponseVerse) string { return v.TextIndopak }
	case "uthmani":
		fnValue = func(v AyahTextResponseVerse) string { return v.TextUthmani }
	case "uthmani_simple":
		fnValue = func(v AyahTextResponseVerse) string { return v.TextUthmaniSimple }
	case "uthmani_tajweed":
		fnValue = func(v AyahTextResponseVerse) string { return v.TextUthmaniTajweed }
	default:
		return fmt.Errorf("unknown text type: %s", name)
	}

	fnClean := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "\n", " ")
		return s
	}

	// Write metadata
	var sb strings.Builder
	title := strings.ReplaceAll(name, "_", " ")
	title = strings.Title(title)
	sb.WriteString("<!--\n")
	sb.WriteString(fmt.Sprintf("Text type: %s\n", title))
	sb.WriteString(fmt.Sprintf("Source   : %s\n", "Quran.com"))
	sb.WriteString("-->\n\n")

	// Write each text
	for i, verse := range ayahText.Verses {
		sb.WriteString("# ")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString("\n\n")

		text := fnClean(fnValue(verse))
		if text == "" {
			sb.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		}

		if i > 0 {
			prevVerse := ayahText.Verses[i-1]
			prevText := fnClean(fnValue(prevVerse))
			if text == prevText {
				sb.WriteString("<!-- TODO:DUPLICATE -->\n\n")
			}
		}

		// Special for tajweed
		if name == "uthmani_tajweed" {
			// Convert text to html.Node document
			doc, err := dom.FastParse(strings.NewReader(text))
			if err != nil {
				return fmt.Errorf("failed to parse uthmani tajweed: %w", err)
			}

			// Replace <tajweed> to []() so it looks better in markdown
			tajweeds := dom.GetElementsByTagName(doc, "tajweed")
			for _, t := range tajweeds {
				tText := dom.TextContent(t)
				tClass := dom.ClassName(t)
				newText := fmt.Sprintf("[%s](%s)", tText, tClass)
				newNode := dom.CreateTextNode(newText)
				dom.ReplaceChild(t.Parent, newNode, t)
			}

			// Remove span
			spans := dom.QuerySelectorAll(doc, "span.end")
			dom.RemoveNodes(spans, nil)

			// Return back doc to text
			body := dom.QuerySelector(doc, "body")
			text = dom.InnerHTML(body)
		} else {
			text = util.MarkdownText(text)
		}

		sb.WriteString(text)
		sb.WriteString("\n\n")
	}

	// Write to file
	dstPath := fmt.Sprintf("%s-qurancom.md", strcase.ToKebab(name))
	dstPath = filepath.Join(dstDir, "ayah-text", dstPath)

	os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	err = os.WriteFile(dstPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write ayah text %q: %w", name, err)
	}

	return nil
}
