package luhn

import (
	"strings"
	"unicode"
)

// Validate checks whether the stringed number is valid according to the Luhn algorithm.
func Validate(number string) bool {
	cleanStr := strings.ReplaceAll(number, " ", "")
	cleanStr = strings.ReplaceAll(cleanStr, "-", "")

	if len(cleanStr) == 0 {
		return false
	}

	for _, r := range cleanStr {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	if len(cleanStr) < 2 {
		return false
	}

	digits := make([]int, len(cleanStr))
	for i, r := range cleanStr {
		digits[i] = int(r - '0')
	}

	for i := len(digits) - 2; i >= 0; i -= 2 {
		digits[i] *= 2
		if digits[i] > 9 {
			digits[i] -= 9
		}
	}

	total := 0
	for _, digit := range digits {
		total += digit
	}

	return total%10 == 0
}
