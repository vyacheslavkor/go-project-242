package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateSize(t *testing.T) {
	testdata := "testdata"

	symlinkRoot := filepath.Join(t.TempDir(), "link")
	require.NoError(t, os.Symlink("in", symlinkRoot))

	testCases := []struct {
		name      string
		path      string
		recursive bool
		all       bool
		expected  int64
	}{
		{
			name:     "file",
			path:     filepath.Join(testdata, "file1.txt"),
			expected: 2906,
		},
		{
			name:     "symlink file",
			path:     symlinkRoot,
			expected: 2,
		},
		{
			name:     "directory",
			path:     filepath.Join(testdata, "in"),
			expected: 2414,
		},
		{
			name:     "directory with hidden",
			path:     filepath.Join(testdata, "in"),
			all:      true,
			expected: 2515,
		},
		{
			name:     "empty directory",
			path:     filepath.Join(testdata, "empty"),
			expected: 0,
		},
		{
			name:      "recursive directory",
			path:      testdata,
			recursive: true,
			expected:  5320,
		},
		{
			name:     "direct hidden file counted with all",
			path:     filepath.Join(testdata, "in", ".file3.txt"),
			all:      true,
			expected: 101,
		},
		{
			name:     "direct hidden directory ignored without all",
			path:     filepath.Join(testdata, ".hidden"),
			expected: 101,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := CalculateSize(
				testCase.path,
				testCase.recursive,
				testCase.all,
			)
			require.NoError(t, err)
			require.Equal(t, testCase.expected, result)
		})
	}

	t.Run("non-existent path", func(t *testing.T) {
		_, err := CalculateSize(filepath.Join(testdata, "missing"), false, false)
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestGetFileSize_FSErrors(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentPath := filepath.Join(tempDir, "missing")

	notPermittedPath := filepath.Join(tempDir, "unreadable")
	require.NoError(t, os.Mkdir(notPermittedPath, 0o000))

	testCases := []struct {
		name      string
		path      string
		recursive bool
		expected  error
		errMsg    string
	}{
		{
			name:      "non-existent path",
			path:      nonExistentPath,
			recursive: false,
			expected:  os.ErrNotExist,
			errMsg:    "failed to read path metadata for " + nonExistentPath,
		},
		{
			name:      "not permitted path",
			path:      notPermittedPath,
			recursive: false,
			expected:  os.ErrPermission,
			errMsg:    "failed to read directory for " + notPermittedPath,
		},
		{
			name:      "not permitted path recursive",
			path:      tempDir,
			recursive: true,
			expected:  os.ErrPermission,
			errMsg:    "failed to read directory for " + tempDir,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := CalculateSize(testCase.path, testCase.recursive, false)
			assert.ErrorIs(t, err, testCase.expected)
			assert.ErrorContains(t, err, testCase.errMsg)
		})
	}
}
