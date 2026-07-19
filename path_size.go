package code

import (
	"code/internal/pathsize"
	"code/internal/sizefmt"
)

// Calculate calculates the size of path and returns it in formatted form.
// Symlinks are measured by link size. Hidden entries nested during directory
// traversal contribute 0 bytes unless all is true. A hidden path passed
// directly as path is always evaluated. Special files also contribute 0 bytes.
func Calculate(path string, recursive, human, all bool) (string, error) {
	size, err := pathsize.Calculate(path, recursive, all)
	if err != nil {
		return "", err
	}

	return sizefmt.Format(size, human), nil
}
