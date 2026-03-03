package utils

// ExtractPhoneDigits strips all non-digit characters from a phone number string,
// preserving the original digit order. This normalizes inputs like:
//
//	"+7 (995) 234 10 96"  → "79952341096"
//	"+998-50-123-45-67"   → "998501234567"
//	"79952341096"         → "79952341096"
func ExtractPhoneDigits(input string) string {
	digits := make([]byte, 0, len(input))
	for i := 0; i < len(input); i++ {
		if input[i] >= '0' && input[i] <= '9' {
			digits = append(digits, input[i])
		}
	}
	return string(digits)
}
