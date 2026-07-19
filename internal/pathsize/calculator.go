package pathsize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Calculate returns the total size in bytes for path.
// Symlinks are measured by link size, not target. Hidden entries nested during
// directory traversal contribute 0 bytes unless all is true. A hidden path
// passed directly as path is always evaluated. Special files are ignored and
// contribute 0 bytes.
func Calculate(path string, recursive, all bool) (int64, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read path metadata for %s: %w", path, err)
	}

	if fileInfo.IsDir() {
		return getDirectorySize(path, recursive, all)
	}

	size, ok := getFileSize(fileInfo)
	if ok {
		return size, nil
	}

	return 0, nil
}

func getDirectorySize(path string, recursive, all bool) (int64, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read directory for %s: %w", path, err)
	}

	result := int64(0)

	for _, file := range files {
		if !all && isHiddenFile(file.Name()) {
			continue
		}

		entrySize, err := processDirEntry(file, path, recursive, all)
		if err != nil {
			return 0, err
		}
		result += entrySize
	}

	return result, nil
}

func getFileSize(fileInfo os.FileInfo) (int64, bool) {
	mode := fileInfo.Mode()
	if !mode.IsRegular() && mode.Type() != os.ModeSymlink {
		return 0, false
	}

	return fileInfo.Size(), true
}

func processDirEntry(file os.DirEntry, path string, recursive, all bool) (int64, error) {
	fileInfo, err := readDirEntryInfo(file, path)
	if err != nil {
		return 0, err
	}

	if fileInfo.IsDir() {
		if !recursive {
			return 0, nil
		}

		return getDirectorySize(
			filepath.Join(path, file.Name()),
			recursive,
			all,
		)
	}

	size, ok := getFileSize(fileInfo)
	if !ok {
		return 0, nil
	}

	return size, nil
}

func readDirEntryInfo(file os.DirEntry, parentPath string) (os.FileInfo, error) {
	fileInfo, err := file.Info()
	if err != nil {
		fullPath := filepath.Join(parentPath, file.Name())
		return nil, fmt.Errorf("failed to read file info for %s: %w", fullPath, err)
	}

	return fileInfo, nil
}

func isHiddenFile(fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}
