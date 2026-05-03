package slugify

import (
	"github.com/essentialkaos/translit/v3"
	"regexp"
	"strings"
)

func Slugify(text string) string {
	transliterator := translit.ICAO
	transliterated := transliterator(text)

	lower := strings.ToLower(transliterated)

	withHyphens := strings.ReplaceAll(lower, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	cleaned := reg.ReplaceAllString(withHyphens, "")

	final := regexp.MustCompile(`-+`).ReplaceAllString(cleaned, "-")
	final = strings.Trim(final, "-")

	return final
}
