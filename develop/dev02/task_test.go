package main

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "valid",
			input:    `a4bc2d5e`,
			expected: `aaaabccddddde`,
			err:      nil,
		},
		{
			name:     "valid",
			input:    `abcd`,
			expected: `abcd`,
			err:      nil,
		},
		{
			name:     "valid",
			input:    ``,
			expected: ``,
			err:      nil,
		},
		{
			name:     "valid",
			input:    `qwe\4\5`,
			expected: `qwe45`,
			err:      nil,
		},
		{
			name:     "valid",
			input:    `qwe\45`,
			expected: `qwe44444`,
			err:      nil,
		},
		{
			name:     "valid",
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
			err:      nil,
		},
		{
			name:     "invalid",
			input:    `45`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			name:     "invalid",
			input:    `\`,
			expected: ``,
			err:      ErrInvalidString,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Unpack(test.input)

			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}

			if err != test.err {
				t.Errorf("Expected %s, got %s", test.err.Error(), err.Error())
			}
		})
	}
}
