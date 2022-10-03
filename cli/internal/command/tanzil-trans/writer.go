package tanzilTrans

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func cleanDstDir(dstDir string) error {
	return filepath.WalkDir(dstDir, func(path string, d fs.DirEntry, err error) error {
		// Remove all file suffixed with "-tanzil.md"
		// that located in of the three directories
		dName := d.Name()
		if d.IsDir() || !strings.HasSuffix(dName, "-tanzil.md") {
			return nil
		}

		dDir := filepath.Base(filepath.Dir(path))
		switch dDir {
		case "ayah-tafsir",
			"ayah-translation",
			"ayah-transliteration":
			return os.Remove(path)
		}

		return nil
	})
}

func writeTranslations(dstDir string, dataList []TranslationData) error {
	for _, data := range dataList {
		logrus.Printf("writing %s", data.FileName)
		if err := writeTranslation(dstDir, data); err != nil {
			return err
		}
	}

	return nil
}

func writeTranslation(dstDir string, data TranslationData) error {
	// Prepare destination path
	dstDir = filepath.Join(dstDir, dataDestination(data.FileName))
	os.MkdirAll(dstDir, os.ModePerm)

	dstPath := data.FileName + "-tanzil.md"
	dstPath = filepath.Join(dstDir, dstPath)

	// Open destination file
	f, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create file for %q failed: %w", data.FileName, err)
	}
	defer f.Close()

	// Write metadata
	f.WriteString("<!--\n")
	f.WriteString(data.Metadata + "\n")
	f.WriteString("-->\n\n")

	// Write each ayah
	for i, ayah := range data.AyahList {
		f.WriteString("# ")
		f.WriteString(strconv.Itoa(i + 1))
		f.WriteString("\n\n")

		if ayah.Empty {
			f.WriteString("<!-- TODO:MISSING -->\n\n")
		} else if ayah.Duplicated {
			f.WriteString("<!-- TODO:DUPLICATE -->\n\n")
		}

		if ayah.Translation != "" {
			f.WriteString(ayah.Translation)
			f.WriteString("\n\n")
		}
	}

	// Flush the data
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("write file for %q failed: %w", data.FileName, err)
	}

	return nil
}

func dataDestination(fileName string) string {
	switch fileName {
	case "ar-jalalayn",
		"ar-muyassar",
		"fa-khorramdel",
		"id-jalalayn",
		"id-muntakhab",
		"ru-kuliev-alsaadi",
		"ru-muntahab",
		"uz-sodik",
		"tr-diyanet":
		return "ayah-tafsir"

	case "en-transliteration",
		"tr-transliteration":
		return "ayah-transliteration"

	default:
		return "ayah-translation"
	}
}
