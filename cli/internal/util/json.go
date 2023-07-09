package util

import (
	"data-quran-cli/internal/norm"
	"encoding/json"
	"os"
	"regexp"
)

var rxJsonKey = regexp.MustCompile(`(?m)^(\s*)"0+(\d+)"(\s*):`)

func EncodeSortedKeyJson(dstPath string, data any) error {
	// Encode JSON
	bt, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Clean up keys
	bt = rxJsonKey.ReplaceAll(bt, []byte("$1\"$2\"$3:"))

	// Save to file
	return os.WriteFile(dstPath, bt, os.ModePerm)
}

func DecodeJsonFile(path string, dst any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := norm.NormalizeReader(f)
	return json.NewDecoder(r).Decode(dst)
}
