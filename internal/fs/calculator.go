package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
	fileInfo, err := file.Info()
	if err != nil {
		fullPath := filepath.Join(path, file.Name())
		return 0, fmt.Errorf("failed to read file info for %s: %w", fullPath, err)
	}

	if fileInfo.IsDir() {
		if !recursive {
			return 0, nil
		}

		subResult, subError := getDirectorySize(
			filepath.Join(path, file.Name()),
			recursive,
			all,
		)
		if subError != nil {
			return 0, subError
		}

		return subResult, nil
	}

	size, ok := getFileSize(fileInfo)
	if !ok {
		return 0, nil
	}

	return size, nil
}

func isHiddenFile(fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}

// CalculateSize returns the total size in bytes for path.
// Symlinks are measured by link size, not target. Hidden files and directories
// contribute 0 bytes unless all is true, including when the hidden path itself
// is passed directly. Special files are ignored and contribute 0 bytes.
func CalculateSize(path string, recursive, all bool) (int64, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read path metadata for %s: %w", path, err)
	}

	result := int64(0)

	if fileInfo.IsDir() {
		dirResult, err := getDirectorySize(path, recursive, all)
		if err != nil {
			return 0, err
		}

		result = dirResult
	} else {
		size, ok := getFileSize(fileInfo)
		if ok {
			result = size
		}
	}

	return result, nil
}
