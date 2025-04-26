package utils

import (
	"html"
	"strings"
)

// SanitizeInput экранирует специальные HTML символы
func SanitizeInput(input string) string {
	return strings.TrimSpace(html.EscapeString(input))
}
