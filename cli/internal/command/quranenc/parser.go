package quranenc

import (
	"data-quran-cli/internal/norm"
	"data-quran-cli/internal/util"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"
)

var rxTransNewLines = regexp.MustCompile(`\s*\n+\s*`)

type Metadata struct {
	Title        string `xml:"title"`
	Language     string `xml:"language"`
	ID           string `xml:"id"`
	Source       string `xml:"source"`
	URL          string `xml:"url"`
	UpdatedAt    string `xml:"updated_at"`
	CheckUpdates string `xml:"check_updates"`
}

type Ayah struct {
	Number      int    `xml:"number,attr"`
	Translation string `xml:"translation"`
	Footnotes   string `xml:"footnotes"`
	Duplicated  bool
}

type Surah struct {
	Number   int    `xml:"number,attr"`
	AyahList []Ayah `xml:"aya"`
}

type TranslationData struct {
	Meta      Metadata `xml:"meta"`
	SurahList []Surah  `xml:"sura_list>sura"`
}

type FlattenedData struct {
	Meta     Metadata
	AyahList []Ayah
}

func parse(cacheDir string) ([]FlattenedData, error) {
	// Get list of file in cache dir
	dirItems, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("read dir failed: %w", err)
	}

	var files []string
	for _, item := range dirItems {
		name := item.Name()
		if !item.IsDir() && filepath.Ext(name) == ".xml" {
			files = append(files, name)
		}
	}

	// Parse each file
	var dataList []FlattenedData
	for _, f := range files {
		logrus.Printf("parsing %s", f)

		fPath := filepath.Join(cacheDir, f)
		data, hasFootnote, err := parseFile(fPath)
		if err != nil {
			return nil, fmt.Errorf("parse %q failed: %w", f, err)
		}

		if hasFootnote {
			logrus.Warnf("%s has footnote", f)
		}

		dataList = append(dataList, *data)
	}

	return dataList, nil
}

func parseFile(path string) (*FlattenedData, bool, error) {
	// Open file
	f, err := os.Open(path)
	if err != nil {
		return nil, false, fmt.Errorf("open file failed: %w", err)
	}
	defer f.Close()

	// Decode XML
	var data TranslationData
	err = xml.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, false, fmt.Errorf("xml decode failed: %w", err)
	}

	// Flatten the data
	var ayahID int
	var hasFootnote bool

	flatData := FlattenedData{Meta: data.Meta}
	for _, surah := range data.SurahList {
		for _, ayah := range surah.AyahList {
			// Normalize text
			translation := norm.NormalizeUnicode(ayah.Translation)
			footnotes := norm.NormalizeUnicode(ayah.Footnotes)

			// Remove newlines in translation
			translation = rxTransNewLines.ReplaceAllString(translation, " ")

			// Make it safe for markdown
			translation = util.MarkdownText(translation)
			footnotes = util.MarkdownText(footnotes)

			// Check if it's duplicated from previous ayah
			var duplicated bool
			if ayahID > 0 {
				prevAyah := flatData.AyahList[ayahID-1]
				prevStr := prevAyah.Translation + "\n" + prevAyah.Footnotes
				str := translation + "\n" + footnotes
				if str == prevStr {
					duplicated = true
				}
			}

			// Check if it has footnote
			hasFootnote = hasFootnote || footnotes != ""

			// Save text
			ayahID++
			flatData.AyahList = append(flatData.AyahList, Ayah{
				Number:      ayahID,
				Translation: translation,
				Footnotes:   footnotes,
				Duplicated:  duplicated,
			})
		}
	}

	// Make sure there are 6236 ayah
	if nAyah := len(flatData.AyahList); nAyah != 6236 {
		return nil, false, fmt.Errorf("n ayah %d != 6236", nAyah)
	}

	return &flatData, hasFootnote, nil
}
