package quranenc

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	iso6391 "github.com/emvi/iso-639-1"
	"github.com/sirupsen/logrus"
)

var (
	rxID = regexp.MustCompile(`(?i)^([^-]+)-(.+)$`)
)

func cleanDstDir(dstDir string) error {
	return filepath.WalkDir(dstDir, func(path string, d fs.DirEntry, err error) error {
		// Remove all file suffixed with "-quranenc.md"
		dName := d.Name()
		if d.IsDir() || !strings.HasSuffix(dName, "-quranenc.md") {
			return nil
		}

		dDir := filepath.Base(filepath.Dir(path))
		if dDir == "ayah-translation" {
			return os.Remove(path)
		}

		return nil
	})
}

func write(dstDir string, dataList []FlattenedData) error {
	for _, data := range dataList {
		logrus.Printf("writing %s", data.Meta.ID)
		if err := writeData(dstDir, data); err != nil {
			return err
		}
	}

	return nil
}

func writeData(dstDir string, data FlattenedData) error {
	// Prepare destination path
	dstDir = filepath.Join(dstDir, "ayah-translation")
	os.MkdirAll(dstDir, os.ModePerm)

	meta := data.Meta
	dstPath := createFileName(meta.ID)
	dstPath = filepath.Join(dstDir, dstPath)

	// Open destination file
	f, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create file for %q failed: %w", meta.ID, err)
	}
	defer f.Close()

	// Write metadata
	f.WriteString("<!--\n")
	f.WriteString(fmt.Sprintf("Title       : %s\n", meta.Title))
	f.WriteString(fmt.Sprintf("Language    : %s\n", meta.Language))
	f.WriteString(fmt.Sprintf("ID          : %s\n", meta.ID))
	f.WriteString(fmt.Sprintf("Source      : %s\n", meta.Source))
	f.WriteString(fmt.Sprintf("URL         : %s\n", meta.URL))
	f.WriteString(fmt.Sprintf("UpdatedAt   : %s\n", meta.UpdatedAt))
	f.WriteString(fmt.Sprintf("CheckUpdates: %s\n", meta.CheckUpdates))
	f.WriteString("-->\n\n")

	// Write each ayah
	for _, ayah := range data.AyahList {
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
		return fmt.Errorf("write file for %q failed: %w", meta.ID, err)
	}

	return nil
}

func createFileName(id string) string {
	// Split ID to several parts
	id = strings.ToLower(id)
	id = strings.ReplaceAll(id, "_", "-")
	parts := rxID.FindStringSubmatch(id)
	if len(parts) != 3 {
		return ""
	}

	// Normalize language
	lang := parts[1]
	lang = upperFirst(lang)
	switch lang {
	case "Ankobambara":
		lang = "Bambara"
	case "Sinhalese":
		lang = "Sinhala"
	case "Azeri":
		lang = "Azerbaijani"
	case "Punjabi":
		lang = "Panjabi"
	}

	// Get language code
	code := iso6391.CodeForName(lang)
	if code == "" {
		return ""
	}

	// Create final name
	return code + "-" + parts[2] + "-quranenc.md"
}

func upperFirst(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
