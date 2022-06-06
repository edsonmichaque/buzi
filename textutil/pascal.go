package textutil

func ToPascalCase(s string) string {
	return toCamelCase(s, true)
}
