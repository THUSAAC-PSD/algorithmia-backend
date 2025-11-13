package constant

import (
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "BodyLimit is correct",
			constant: BodyLimit,
			expected: "2M",
		},
		{
			name:     "GzipLevel is correct",
			constant: string(rune(GzipLevel + '0')),
			expected: "5",
		},
		{
			name:     "Dev environment constant",
			constant: Dev,
			expected: "development",
		},
		{
			name:     "Test environment constant",
			constant: Test,
			expected: "test",
		},
		{
			name:     "Production environment constant",
			constant: Production,
			expected: "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.constant)
			}
		})
	}
}
