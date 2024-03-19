package helper

import "unicode"

func CamelToCapitalize(word string) string {
	var result []rune

	for i, char := range word {
		if i > 0 && unicode.IsUpper(char) {
			result = append(result, ' ')
		}
		result = append(result, char)
	}

	return string(result)
}
