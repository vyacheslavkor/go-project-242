package code

import (
	"code/internal/fs"
	"path/filepath"
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
			name:     "formats size",
			path:     filepath.Join("testdata", "file1.txt"),
			expected: "2906B",
		},
		{
			name:     "human readable output",
			path:     filepath.Join("testdata", "in"),
			human:    true,
			expected: "2.4KB",
		},
		{
			name:     "direct hidden file ignored without all",
			path:     filepath.Join("testdata", "in", ".file3.txt"),
			expected: "0B",
		},
		{
			name:     "direct hidden file counted with all",
			path:     filepath.Join("testdata", "in", ".file3.txt"),
			all:      true,
			expected: "100B",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := GetPathSize(
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
		_, err := GetPathSize(filepath.Join("testdata", "missing"), false, false, false)
		require.ErrorIs(t, err, fs.ErrPathNotExist)
	})
}
