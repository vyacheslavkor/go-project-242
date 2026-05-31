package code_test

import (
	"code"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatSize(t *testing.T) {
	testCases := []struct {
		name     string
		size     int64
		human    bool
		expected string
	}{
		{"zero bytes raw", 0, false, "0B"},
		{"large bytes raw", 1048576, false, "1048576B"},
		{"zero bytes human", 0, true, "0B"},
		{"max bytes before KB", 1023, true, "1023B"},
		{"exact 1 KB boundary", 1024, true, "1.0KB"},
		{"fractional KB round down", 1100, true, "1.1KB"},
		{"fractional KB round up", 1500, true, "1.5KB"},
		{"exact 1 MB", 1048576, true, "1.0MB"},
		{"fractional MB", 1572864, true, "1.5MB"},
		{"exact 1 GB", 1073741824, true, "1.0GB"},
		{"exact 1 TB", 1099511627776, true, "1.0TB"},
		{"exact 1 PB", 1125899906842624, true, "1.0PB"},
		{"exact 1 EB", 1152921504606846976, true, "1.0EB"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := code.FormatSize(
				testCase.size,
				testCase.human,
			)
			require.Equal(t, testCase.expected, result)
		})
	}
}
