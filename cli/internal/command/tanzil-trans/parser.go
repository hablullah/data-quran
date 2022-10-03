package tanzilTrans

import (
	"bufio"
	"data-quran-cli/internal/norm"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Ayah struct {
	Translation string
	Duplicated  bool
	Empty       bool
}

type TranslationData struct {
	FileName string
	Metadata string
	AyahList []Ayah
}

var (
	rxMetaSeparator = regexp.MustCompile(`^-{3,}$`)
)

func parse(cacheDir string) ([]TranslationData, error) {
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
	var dataList []TranslationData
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

func parseFile(path string) (*TranslationData, error) {
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
		// Fetch the line
		line := scanner.Text()
		line = strings.TrimPrefix(line, "#")
		line = strings.TrimSpace(line)

		// Normalize the line
		line = norm.NormalizeUnicode(line)

		// Check if metadata started
		if rxMetaSeparator.MatchString(line) {
			metadataStarted = true
			continue
		}

		// Put the line in its respective place
		if !metadataStarted && ayahIdx < nAyah {
			ayahList[ayahIdx] = Ayah{Translation: line}
			ayahIdx++
		} else {
			metadata += line + "\n"
		}
	}

	// Mark missing or duplicate lines
	for i, ayah := range ayahList {
		// Check if translation is empty
		// #NÁZEV? is used specifically for Czech language
		ayah.Empty = ayah.Translation == "" || ayah.Translation == "#NÁZEV?"

		// Check if it's duplicate of previous ayah
		if i > 0 && ayah.Translation == ayahList[i-1].Translation {
			ayah.Duplicated = true
		}

		// Special: ayah 6096 (Asy-Syarh 94:6) and 6172 (Ath-Takathur 102:4)
		// is not duplicate because it's actually has same translation as
		// the previous ayah.
		switch i + 1 {
		case 6096, 6172:
			ayah.Duplicated = false
		}

		ayahList[i] = ayah
	}

	// Return the final data
	fileName := filepath.Base(path)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	return &TranslationData{
		FileName: fileName,
		AyahList: ayahList,
		Metadata: strings.TrimSpace(metadata),
	}, nil
}
