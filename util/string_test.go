package util

import (
	"testing"
	"unicode"
)

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		expected int
	}{
		{
			name:     "Generate string of length 0",
			length:   0,
			expected: 0,
		},
		{
			name:     "Generate string of length 1",
			length:   1,
			expected: 1,
		},
		{
			name:     "Generate string of length 10",
			length:   10,
			expected: 10,
		},
		{
			name:     "Generate string of length 50",
			length:   50,
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateRandomString(tt.length)

			// Check if the length of the generated string matches the expected length
			if len(result) != tt.expected {
				t.Errorf("Expected length %d, got %d", tt.expected, len(result))
			}

			// Check if the string contains only valid characters
			for _, char := range result {
				if !(unicode.IsLetter(char) || unicode.IsDigit(char)) {
					t.Errorf("Generated string contains invalid character: %c", char)
				}
			}
		})
	}
}
