package helpers

import (
	"strings"
	"unicode"
)

// Titleize ensure the provided string has uppercase initials
func Titleize(input string) (titleized string) {
	isToUpper := false
	for k, v := range input {
		if k == 0 {
			titleized = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper || unicode.IsUpper(v) {
				titleized += " " + strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if (v == '_') || (v == ' ') {
					isToUpper = true
				} else {
					titleized += string(v)
				}
			}
		}
	}
	return

}
