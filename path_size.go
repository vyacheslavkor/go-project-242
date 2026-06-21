package code

import (
	"code/internal/formatter"
	"code/internal/fs"
)

// GetPathSize calculates the size of path and returns it in formatted form.
// Symlinks are measured by link size. Hidden files and directories contribute
// 0 bytes unless all is true, including when the hidden path itself is passed
// directly. Special files also contribute 0 bytes.
func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := fs.CalculateSize(path, recursive, all)
	if err != nil {
		return "", err
	}

	return formatter.FormatSize(size, human), nil
}
