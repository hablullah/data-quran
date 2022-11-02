package util

import (
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

var mdc = md.NewConverter("", true, nil)

func MarkdownText(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	s, err := mdc.ConvertString(s)
	if err != nil {
		panic(err)
	}

	return s
}
