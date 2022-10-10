package kemenag

import (
	"data-quran-cli/internal/norm"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-shiori/dom"
	"github.com/sirupsen/logrus"
)

type Tafsir struct {
	TafsirID      int    `json:"tafsir_id"`
	SurahID       int    `json:"surah_id"`
	AyaNumber     int    `json:"aya_number"`
	TafsirWajiz   string `json:"tafsir_wajiz"`
	TafsirTahlili string `json:"tafsir_tahlili"`
}

type TafsirContainer struct {
	Msg    string   `json:"msg"`
	Tafsir []Tafsir `json:"tafsir"`
}

func parseAllTafsir(cacheDir string) ([]Tafsir, error) {
	// Get list of file in cache dir
	dirItems, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("read dir for tafsir failed: %w", err)
	}

	var files []string
	for _, item := range dirItems {
		name := item.Name()
		ext := filepath.Ext(name)
		if !item.IsDir() && strings.HasPrefix(name, "ayah-") && ext == ".json" {
			files = append(files, name)
		}
	}

	// Parse each file
	var tafsirs []Tafsir
	for _, f := range files {
		logrus.Printf("parsing %s", f)
		srcPath := filepath.Join(cacheDir, f)
		currentTafsirs, err := parseTafsir(srcPath)
		if err != nil {
			return nil, err
		}

		tafsirs = append(tafsirs, currentTafsirs...)
	}

	// Make sure there is 6236 ayah
	if nTafsir := len(tafsirs); nTafsir != 6236 {
		return nil, fmt.Errorf("n tafsir %d != 6236", nTafsir)
	}

	return tafsirs, nil
}

func parseTafsir(srcPath string) ([]Tafsir, error) {
	// Open and decode source file
	srcName := filepath.Base(srcPath)
	src, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", srcName, err)
	}
	defer src.Close()

	var srcData TafsirContainer
	err = json.NewDecoder(src).Decode(&srcData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", srcName, err)
	}

	// Normalize data
	tafsirs := append([]Tafsir{}, srcData.Tafsir...)
	for i, tafsir := range tafsirs {
		// Normalize unicode
		wajiz := norm.NormalizeUnicode(tafsir.TafsirWajiz)
		tahlili := norm.NormalizeUnicode(tafsir.TafsirTahlili)

		// Clean up
		wajiz = cleanTafsirWajiz(wajiz)
		tahlili = cleanTafsirTahlili(tahlili)

		tafsir.TafsirWajiz = wajiz
		tafsir.TafsirTahlili = tahlili
		tafsirs[i] = tafsir
	}

	return tafsirs, nil
}

func cleanTafsirWajiz(s string) string {
	// Remove html tags
	div := dom.CreateElement("div")
	dom.SetInnerHTML(div, s)
	s = dom.TextContent(div)

	// Remove weird conversion
	s = strings.ReplaceAll(s, "no_ayat", "aya")
	s = strings.ReplaceAll(s, "no_surah", "sura")
	s = strings.TrimSpace(s)

	if strings.Contains(s, "\n") {
		return ""
	}

	return s
}

func cleanTafsirTahlili(s string) string {
	// Replace semicolon
	s = strings.ReplaceAll(s, ";", ".\n")

	// Skip blank line
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n\n")
}
