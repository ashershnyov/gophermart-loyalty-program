package luhn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		number   string
		expected bool
	}{
		{"79927398713", true},
		{"4242424242424242", true},
		{"1234567890", false},
		{"1", false},
		{"12", false},
		{"123", false},
		{"4532015112830366", true},
		{"6011514433546201", true},
		{"6011514433546202", false},
		{"79927398710", false},
		{"7992 7398", true},
		{"7992-7398", true},
		{"79ab-7398", false},
		{"---", false},
		{"", false},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.expected, Validate(tt.number))
		})
	}
}
