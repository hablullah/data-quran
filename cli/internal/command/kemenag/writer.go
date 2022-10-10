package kemenag

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/zyedidia/generic/mapset"
)

func writeTafsir(dstDir, name string, tafsirs []string) error {
	logrus.Printf("writing tafsir %s", name)

	// Prepare destination path
	dstDir = filepath.Join(dstDir, "ayah-tafsir")
	os.MkdirAll(dstDir, os.ModePerm)

	// Open destination file
	dstPath := filepath.Join(dstDir, name+".md")
	f, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create file for tafsir %s failed: %w", name, err)
	}
	defer f.Close()

	// Check for duplicate
	duplicateIdxs := mapset.New[int]()
	for i := 1; i < len(tafsirs); i++ {
		if tafsirs[i] == tafsirs[i-1] {
			duplicateIdxs.Put(i)
		}
	}

	// Write each tafsir
	for i, tafsir := range tafsirs {
		f.WriteString("# ")
		f.WriteString(strconv.Itoa(i + 1))
		f.WriteString("\n\n")

		if tafsir == "" {
			f.WriteString("<!-- TODO:MISSING -->\n\n")
			continue
		}

		if duplicateIdxs.Has(i) {
			f.WriteString("<!-- TODO:DUPLICATE -->\n\n")
		}

		f.WriteString(tafsir)
		f.WriteString("\n\n")
	}

	// Flush the data
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("write file for tafsir %s failed: %w", name, err)
	}

	return nil
}
