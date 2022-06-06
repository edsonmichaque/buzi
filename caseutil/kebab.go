package textutil

import (
	"bytes"
	"unicode"
)

func ToKebab(s string) string {
	var buf bytes.Buffer

	for i, c := range s {
		if buf.Len() == 0 {
			if !unicode.IsLetter(c) {
				continue
			}

			buf.WriteRune(unicode.ToLower(c))
		} else if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			previousRune := rune(s[i-1])
			if unicode.IsLetter(previousRune) || unicode.IsNumber(previousRune) {
				buf.WriteRune('-')
			}
		} else if unicode.IsLetter(c) {
			if unicode.IsUpper(c) {
				buf.WriteRune(unicode.ToLower('-'))
			}

			buf.WriteRune(unicode.ToLower(c))
		} else {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}
