package tanzilText

import (
	"bufio"
	"data-quran-cli/internal/norm"
	"data-quran-cli/internal/util"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/agext/levenshtein"
	"github.com/sirupsen/logrus"
	"github.com/zyedidia/generic/mapset"
)

type Ayah struct {
	Text       string
	Duplicated bool
	Empty      bool
}

type TextData struct {
	FileName string
	Metadata string
	AyahList []Ayah
}

var (
	rxMetaPrefix    = regexp.MustCompile(`^#\s{0,2}`)
	rxMetaSeparator = regexp.MustCompile(`^={3,}$`)
)

func parse(cacheDir string) ([]TextData, error) {
	// Get list of file in cache dir
	dirItems, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("read dir failed: %w", err)
	}

	var files []string
	for _, item := range dirItems {
		name := item.Name()
		if !item.IsDir() && filepath.Ext(name) == ".txt" {
			files = append(files, name)
		}
	}

	// Parse each file
	var dataList []TextData
	for _, f := range files {
		logrus.Printf("parsing %s", f)

		fPath := filepath.Join(cacheDir, f)
		data, err := parseFile(fPath)
		if err != nil {
			return nil, fmt.Errorf("parse %q failed: %w", f, err)
		}

		dataList = append(dataList, *data)
	}

	return dataList, nil
}

func parseFile(path string) (*TextData, error) {
	// Open file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer f.Close()

	// Decode texts
	nAyah := 6236
	ayahList := make([]Ayah, nAyah)
	var ayahIdx int
	var metadata string
	var metadataStarted bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Fetch and normalize the line
		line := scanner.Text()
		line = norm.NormalizeUnicode(line)

		// Check if metadata started
		if !metadataStarted && strings.HasPrefix(line, "#") {
			metadataStarted = true
		}

		// Put the line in its respective place
		if !metadataStarted && ayahIdx < nAyah {
			line = util.MarkdownText(line)
			ayahList[ayahIdx] = Ayah{Text: line}
			ayahIdx++
		} else {
			line = rxMetaPrefix.ReplaceAllString(line, "")
			if !rxMetaSeparator.MatchString(line) {
				metadata += line + "\n"
			}
		}
	}

	// Mark missing or duplicate lines and remove basmalah
	basmalah := ayahList[0].Text
	nBasmalahWords := len(strings.Fields(basmalah))

	for i, ayah := range ayahList {
		// Check if translation is empty
		ayah.Empty = ayah.Text == ""

		// Check if it's duplicate of previous ayah
		if i > 0 && ayah.Text == ayahList[i-1].Text {
			ayah.Duplicated = true
		}

		// Remove basmalah
		if starterAyah.Has(i + 1) {
			text := ayah.Text
			if strings.HasPrefix(text, basmalah) {
				text = strings.TrimPrefix(text, basmalah)
			} else {
				words := strings.Fields(text)
				possiblyBasmalah := strings.Join(words[:nBasmalahWords], " ")
				similarity := levenshtein.Similarity(basmalah, possiblyBasmalah, nil)
				if similarity >= 0.8 {
					text = strings.Join(words[nBasmalahWords:], " ")
				}
			}

			ayah.Text = strings.TrimSpace(text)
		}

		ayahList[i] = ayah
	}

	// Return the final data
	fileName := filepath.Base(path)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	return &TextData{
		FileName: fileName,
		AyahList: ayahList,
		Metadata: strings.TrimSpace(metadata),
	}, nil
}

var starterAyah mapset.Set[int]

func init() {
	listStarterAyah := []int{
		8, 294, 494, 670, 790, 955, 1161, 1236, 1365, 1474, 1597, 1708, 1751,
		1803, 1902, 2030, 2141, 2251, 2349, 2484, 2596, 2674, 2792, 2856, 2933,
		3160, 3253, 3341, 3410, 3470, 3504, 3534, 3607, 3661, 3706, 3789, 3971,
		4059, 4134, 4219, 4273, 4326, 4415, 4474, 4511, 4546, 4584, 4613, 4631,
		4676, 4736, 4785, 4847, 4902, 4980, 5076, 5105, 5127, 5151, 5164, 5178,
		5189, 5200, 5218, 5230, 5242, 5272, 5324, 5376, 5420, 5448, 5476, 5496,
		5552, 5592, 5623, 5673, 5713, 5759, 5801, 5830, 5849, 5885, 5910, 5932,
		5949, 5968, 5994, 6024, 6044, 6059, 6080, 6091, 6099, 6107, 6126, 6131,
		6139, 6147, 6158, 6169, 6177, 6180, 6189, 6194, 6198, 6205, 6208, 6214,
		6217, 6222, 6226, 6231,
	}

	starterAyah = mapset.New[int]()
	for _, l := range listStarterAyah {
		starterAyah.Put(l)
	}
}
