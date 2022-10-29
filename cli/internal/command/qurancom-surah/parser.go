package quranComSurah

import (
	"data-quran-cli/internal/norm"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-shiori/dom"
	"github.com/sirupsen/logrus"
)

type ListSurahSource struct {
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
		ChapterNumber   int    `json:"chapter_number"`
		TranslatedName  struct {
			LanguageName string `json:"language_name"`
			Name         string `json:"name"`
		} `json:"translated_name"`
	} `json:"chapters"`
}

type SurahInfoSource struct {
	ChapterInfo struct {
		ChapterID    int    `json:"chapter_id"`
		LanguageName string `json:"language_name"`
		ShortText    string `json:"short_text"`
		Source       string `json:"source"`
		Text         string `json:"text"`
	} `json:"chapter_info"`
}

type ListSurahOutput struct {
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

type SurahInfoOutput struct {
	Number   int
	Language string
	Source   string
	Text     string
}

type AllSurahInfoOutput struct {
	Language string
	Source   string
	Texts    map[int]string
}

func parseListSurah(cacheDir string, language string) (map[string]ListSurahOutput, error) {
	// Open file
	path := fmt.Sprintf("list-%s.json", language)
	path = filepath.Join(cacheDir, path)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail to open surah list for %s: %w", language, err)
	}
	defer f.Close()

	// Decode data
	var srcData ListSurahSource
	err = json.NewDecoder(f).Decode(&srcData)
	if err != nil {
		return nil, fmt.Errorf("fail to decode surah list for %s: %w", language, err)
	}

	// Generate output
	output := map[string]ListSurahOutput{}
	for _, chapter := range srcData.Chapters {
		// In Quran.com, if data for a language not exist, they will
		// fallback into using English language. In this case, we
		// will just skip it.
		if language != "en" && chapter.TranslatedName.LanguageName == "english" {
			continue
		}

		key := fmt.Sprintf("%03d", chapter.ChapterNumber)
		output[key] = ListSurahOutput{
			Name:        norm.NormalizeUnicode(chapter.NameSimple),
			Translation: norm.NormalizeUnicode(chapter.TranslatedName.Name),
		}
	}

	// Check if translation complete
	if n := len(output); n != 114 {
		logrus.Warnf("surah list for %s: want 114 got %d", language, n)
		return nil, nil
	}

	return output, nil
}

func parseAllSurahInfo(cacheDir, language string, mdc *md.Converter) (*AllSurahInfoOutput, error) {
	// Extract each surah in this language
	mapInfo := map[int]string{}
	var languageName, source string

	for surah := 1; surah <= 114; surah++ {
		output, err := parseSurahInfo(cacheDir, language, surah, mdc)
		if err != nil {
			return nil, err
		} else if output == nil {
			continue
		}

		if languageName == "" || source == "" {
			source = output.Source
			languageName = output.Language
		}

		mapInfo[surah] = output.Text
	}

	// Check if info complete
	if n := len(mapInfo); n != 114 {
		logrus.Warnf("surah info for %s: want 114 got %d", language, n)
		if n == 0 {
			return nil, nil
		}
	}

	return &AllSurahInfoOutput{
		Source:   source,
		Language: languageName,
		Texts:    mapInfo,
	}, nil
}

func parseSurahInfo(cacheDir, language string, surah int, mdc *md.Converter) (*SurahInfoOutput, error) {
	// Open file
	path := fmt.Sprintf("info-%s-%03d.json", language, surah)
	path = filepath.Join(cacheDir, path)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail to open surah info for %s %d: %w", language, surah, err)
	}
	defer f.Close()

	// Decode data
	var srcData SurahInfoSource
	err = json.NewDecoder(f).Decode(&srcData)
	if err != nil {
		return nil, fmt.Errorf("fail to decode surah info for %s %d: %w", language, surah, err)
	}

	// In Quran.com, if data for a language not exist, they will
	// fallback into using English language. In this case, we
	// will just skip it.
	if language != "en" && srcData.ChapterInfo.LanguageName == "english" {
		return nil, nil
	}

	// If text is empty, just stop
	srcText := norm.NormalizeUnicode(srcData.ChapterInfo.Text)
	if srcText == "" {
		return nil, nil
	}

	// Convert text to html.Node document
	doc, err := dom.FastParse(strings.NewReader(srcText))
	if err != nil {
		return nil, fmt.Errorf("fail to parse surah info HTML for %s %d: %w", language, surah, err)
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

	// Convert text to markdown
	markdown, err := mdc.ConvertString(docHTML)
	if err != nil {
		return nil, fmt.Errorf("fail to create surah info md for %s %d: %w", language, surah, err)
	}

	// Return output
	return &SurahInfoOutput{
		Text:     markdown,
		Number:   srcData.ChapterInfo.ChapterID,
		Source:   norm.NormalizeUnicode(srcData.ChapterInfo.Source),
		Language: norm.NormalizeUnicode(srcData.ChapterInfo.LanguageName),
	}, nil
}
