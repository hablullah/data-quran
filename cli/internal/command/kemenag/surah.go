package kemenag

import (
	"data-quran-cli/internal/norm"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/sirupsen/logrus"
)

type Ayah struct {
	IDAyat       int    `json:"id_ayat"`
	NoSurah      int    `json:"no_surah"`
	NoAyat       int    `json:"no_ayat"`
	TeksAyat     string `json:"teks_ayat"`
	Tema         string `json:"tema"`
	TeksTerjemah string `json:"teks_terjemah"`
	NoFn         string `json:"no_fn"`
	TeksFn       string `json:"teks_fn"`
}

type Surah struct {
	Data []Ayah `json:"data"`
}

type AyahOutput struct {
	Number      int
	Translation string
	Footnotes   string
	Duplicated  bool
}

var (
	rxNumberOnly     = regexp.MustCompile(`\d+`)
	rxFootnoteNumber = regexp.MustCompile(`(?m)^(\d+)\)\s*`)
)

func parseAndWriteAllSurah(cacheDir, dstDir string) error {
	listAyah, err := parseAllSurah(cacheDir)
	if err != nil {
		return err
	}

	err = writeAllSurah(dstDir, listAyah)
	return err
}

func parseAllSurah(cacheDir string) ([]AyahOutput, error) {
	// Get list of file in cache dir
	dirItems, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("read dir for surah failed: %w", err)
	}

	var files []string
	for _, item := range dirItems {
		name := item.Name()
		ext := filepath.Ext(name)
		if !item.IsDir() && strings.HasPrefix(name, "surah-") && ext == ".json" {
			files = append(files, name)
		}
	}

	// Parse each file
	var outputs []AyahOutput
	for _, f := range files {
		logrus.Printf("parsing %s", f)
		srcPath := filepath.Join(cacheDir, f)
		surahOutputs, err := parseSurah(srcPath)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, surahOutputs...)
	}

	return outputs, nil
}

func parseSurah(srcPath string) ([]AyahOutput, error) {
	// Open and decode source file
	srcName := filepath.Base(srcPath)
	src, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", srcName, err)
	}
	defer src.Close()

	var srcData Surah
	err = json.NewDecoder(src).Decode(&srcData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", srcName, err)
	}

	// Convert the data
	div := dom.CreateElement("div")
	outputs := make([]AyahOutput, len(srcData.Data))
	for i, d := range srcData.Data {
		// Normalize unicode
		trans := norm.NormalizeUnicode(d.TeksTerjemah)
		footnote := norm.NormalizeUnicode(d.TeksFn)

		// Convert HTML tags in translation
		dom.SetInnerHTML(div, trans)
		for _, sup := range dom.QuerySelectorAll(div, "sup") {
			supText := dom.TextContent(sup)
			fnNumber := rxNumberOnly.FindString(supText)
			fnNode := dom.CreateTextNode("[^" + fnNumber + "]")
			dom.ReplaceChild(div, fnNode, sup)
		}
		trans = dom.TextContent(div)

		// Convert HTML tags in footnote
		dom.SetInnerHTML(div, footnote)
		for _, br := range dom.QuerySelectorAll(div, "br") {
			newlineNode := dom.CreateTextNode("\n\n")
			dom.ReplaceChild(div, newlineNode, br)
		}

		footnote = dom.TextContent(div)
		footnote = rxFootnoteNumber.ReplaceAllString(footnote, "[^$1]: ")

		outputs[i] = AyahOutput{
			Number:      d.IDAyat,
			Translation: strings.TrimSpace(trans),
			Footnotes:   strings.TrimSpace(footnote),
		}
	}

	// Mark for duplicate
	for i := 1; i < len(outputs); i++ {
		current := outputs[i].Translation
		previous := outputs[i-1].Translation
		outputs[i].Duplicated = current == previous
	}

	return outputs, nil
}

func writeAllSurah(dstDir string, listAyah []AyahOutput) error {
	logrus.Println("writing surah translation")

	// Prepare destination path
	dstDir = filepath.Join(dstDir, "ayah-translation")
	os.MkdirAll(dstDir, os.ModePerm)

	// Open destination file
	dstPath := filepath.Join(dstDir, "id-kemenag.md")
	f, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create file for all surah failed: %w", err)
	}
	defer f.Close()

	// Write each ayah
	for _, ayah := range listAyah {
		f.WriteString("# ")
		f.WriteString(strconv.Itoa(ayah.Number))
		f.WriteString("\n\n")

		if ayah.Translation == "" {
			f.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		}

		if ayah.Duplicated {
			f.WriteString("<!-- TODO:DUPLICATE -->\n\n")
		}

		f.WriteString(ayah.Translation)
		f.WriteString("\n\n")

		if ayah.Footnotes != "" {
			f.WriteString(ayah.Footnotes)
			f.WriteString("\n\n")
		}
	}

	// Flush the data
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("write file for all surah failed: %w", err)
	}

	return nil
}
