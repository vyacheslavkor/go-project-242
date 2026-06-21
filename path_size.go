package code

import (
	"code/internal/formatter"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type UserError struct {
	Message string
	Path    string
}

type Namer interface {
	Name() string
}

func (e *UserError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Path)
}

func filterHiddenFiles(files []os.DirEntry) []os.DirEntry {
	result := []os.DirEntry{}
	for _, file := range files {
		if !isHiddenFile(file) {
			result = append(result, file)
		}
	}
	return result
}

func getDirectorySize(path string, recursive, all bool) (int64, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return 0, &UserError{Message: "permission denied", Path: path}
		}

		return 0, fmt.Errorf("failed to read directory: %w", err)
	}

	if !all {
		files = filterHiddenFiles(files)
	}

	result := int64(0)

	for _, file := range files {
		entrySize, err := processDirEntry(file, path, recursive, all)
		if err != nil {
			return 0, err
		}
		result += entrySize
	}

	return result, nil
}

func processDirEntry(file os.DirEntry, path string, recursive, all bool) (int64, error) {
	fileInfo, err := file.Info()
	if err != nil {
		return 0, fmt.Errorf("failed to read file info: %w", err)
	}

	if fileInfo.IsDir() {
		if recursive {
			subResult, subError := getDirectorySize(
				filepath.Join(path, file.Name()),
				recursive,
				all,
			)
			if subError != nil {
				return 0, subError
			}
			return subResult, nil
		} else {
			return 0, nil
		}
	}

	if !fileInfo.Mode().IsRegular() {
		return 0, nil
	}

	return fileInfo.Size(), nil
}

func isHiddenFile(file Namer) bool {
	return strings.HasPrefix(file.Name(), ".")
}

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", &UserError{Message: "path not exists", Path: path}
		}

		return "", fmt.Errorf("failed to read path metadata: %w", err)
	}

	result := int64(0)

	if fileInfo.IsDir() {
		dirResult, err := getDirectorySize(path, recursive, all)
		if err != nil {
			return "", err
		}

		result = dirResult
	} else {
		mode := fileInfo.Mode()
		switch {
		case !all && isHiddenFile(fileInfo):
			result = 0
		case mode.IsRegular(), mode.Type() == os.ModeSymlink:
			result = fileInfo.Size()
		default:
			result = 0
		}
	}

	return formatter.FormatSize(result, human), nil
}
