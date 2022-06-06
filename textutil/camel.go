package textutil

import (
	"bytes"
	"unicode"
)

func ToCamelCase(s string) string {
	return toCamelCase(s, false)
}

func toCamelCase(s string, firstUpperCase bool) string {
	var buf bytes.Buffer

	setUpperCase := false

	for _, c := range s {
		if buf.Len() == 0 {
			if !unicode.IsLetter(c) {
				continue
			}

			if firstUpperCase {
				buf.WriteRune(unicode.ToUpper(c))
			} else {
				buf.WriteRune(unicode.ToLower(c))
			}
		} else if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			setUpperCase = true
		} else if unicode.IsLetter(c) {
			if setUpperCase {
				buf.WriteRune(unicode.ToUpper(c))
				setUpperCase = false
			} else {
				buf.WriteRune(c)
			}
		} else if unicode.IsNumber(c) {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}
