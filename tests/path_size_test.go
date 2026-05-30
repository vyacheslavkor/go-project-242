package code_test

import (
	"code"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathSize(t *testing.T) {
	testCases := []struct {
		name      string
		path      string
		recursive bool
		human     bool
		all       bool
		expected  string
	}{
		{
			name:      "file",
			path:      "./../testdata/file1.txt",
			recursive: false,
			human:     false,
			all:       false,
			expected:  "2906B",
		},
		{
			name:      "directory",
			path:      "./../testdata/in",
			recursive: false,
			human:     false,
			all:       false,
			expected:  "2412B",
		},
		{
			name:      "directory human",
			path:      "./../testdata/in",
			recursive: false,
			human:     true,
			all:       false,
			expected:  "2.4KB",
		},
		{
			name:      "directory hidden",
			path:      "./../testdata/in",
			recursive: false,
			human:     false,
			all:       true,
			expected:  "2512B",
		},
		{
			name:      "empty directory",
			path:      "./../testdata/empty",
			recursive: false,
			human:     false,
			all:       false,
			expected:  "0B",
		},
		{
			name:      "recursive directory",
			path:      "./../testdata",
			recursive: true,
			human:     false,
			all:       false,
			expected:  "5318B",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := code.GetPathSize(
				testCase.path,
				testCase.recursive,
				testCase.human,
				testCase.all,
			)
			require.NoError(t, err)
			require.Equal(t, testCase.expected, result)
		})
	}

	t.Run("non-existent path", func(t *testing.T) {
		_, err := code.GetPathSize("./../testdata/adjfgjkdasf", false, false, false)
		require.Error(t, err)
	})
}
