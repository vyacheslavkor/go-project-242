package code

import (
	"code/internal/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathSize(t *testing.T) {
	t.Run("composes calculate and raw format", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "file.txt")
		require.NoError(t, os.WriteFile(path, []byte("hello"), 0o600))

		result, err := GetPathSize(path, false, false, false)
		require.NoError(t, err)
		require.Equal(t, "5B", result)
	})

	t.Run("composes calculate and human format", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "file.txt")
		content := make([]byte, 1500)
		for i := range content {
			content[i] = 'x'
		}
		require.NoError(t, os.WriteFile(path, content, 0o600))

		result, err := GetPathSize(path, false, true, false)
		require.NoError(t, err)
		require.Equal(t, "1.5KB", result)
	})

	t.Run("propagates path not exist error", func(t *testing.T) {
		_, err := GetPathSize(filepath.Join(t.TempDir(), "missing"), false, false, false)
		require.ErrorIs(t, err, fs.ErrPathNotExist)
	})
}
