package norm

import (
	"strings"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	unicodeNormalizer = transform.Chain(norm.NFKD, norm.NFKC)
)

func NormalizeUnicode(s string) string {
	result, _, err := transform.String(unicodeNormalizer, s)
	if err == nil {
		return strings.TrimSpace(result)
	}
	return s
}
