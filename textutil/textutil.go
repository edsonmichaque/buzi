package textutil

func hasRune(runes []rune, item rune) bool {
	for _, r := range runes {
		if r == item {
			return true
		}
	}

	return false
}

func hasString(list []string, item string) bool {
	for _, s := range list {
		if s == item {
			return true
		}
	}

	return false
}
