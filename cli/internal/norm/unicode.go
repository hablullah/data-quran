package norm

import (
	"strings"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/rangetable"
)

var (
	rtlRuneRemover    = runes.Remove(runes.In(rangetable.New('\u200f')))
	unicodeNormalizer = transform.Chain(norm.NFKD, rtlRuneRemover, norm.NFKC)
)

func NormalizeUnicode(s string) string {

	result, _, err := transform.String(unicodeNormalizer, s)
	if err == nil {
		return strings.TrimSpace(result)
	}
	return s
}
