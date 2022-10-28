package util

import (
	"encoding/json"
	"os"
	"regexp"
)

var rxJsonKey = regexp.MustCompile(`(?m)^(\s*)"0+(\d+)"(\s*):`)

func EncodeSortedKeyJson(dstPath string, data any) error {
	// Encode JSON
	bt, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	// Clean up keys
	bt = rxJsonKey.ReplaceAll(bt, []byte("$1\"$2\"$3:"))

	// Save to file
	return os.WriteFile(dstPath, bt, os.ModePerm)
}
