package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func formatSize(size int64, human bool) string {
	if human {
		return formatToHuman(size)
	}
	return fmt.Sprintf("%dB", size)
}

func formatToHuman(size int64) string {
	sizeFloat := float64(size)
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	unitIndex := 0
	sizeStepMultiplier := float64(1024)
	for sizeFloat >= sizeStepMultiplier && unitIndex < len(units)-1 {
		sizeFloat /= sizeStepMultiplier
		unitIndex++
	}

	if unitIndex == 0 {
		return fmt.Sprintf("%.0f%s", sizeFloat, units[unitIndex])
	}

	return fmt.Sprintf("%.1f%s", sizeFloat, units[unitIndex])
}

func filterHiddenFiles(files []os.DirEntry) []os.DirEntry {
	result := []os.DirEntry{}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			result = append(result, file)
		}
	}
	return result
}

func getDirectorySize(path string, recursive, all bool) (int64, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	if !all {
		files = filterHiddenFiles(files)
	}

	result := int64(0)

	if len(files) == 0 {
		return result, nil
	}

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return 0, err
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
				result += subResult
			} else {
				continue
			}
		} else {
			if !fileInfo.Mode().IsRegular() {
				continue
			}

			result += fileInfo.Size()
		}
	}

	return result, nil
}

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	result := int64(0)

	if fileInfo.IsDir() {
		dirResult, err := getDirectorySize(path, recursive, all)
		if err != nil {
			return "", err
		}

		result = dirResult
	} else {
		switch {
		case !fileInfo.Mode().IsRegular(), !all && strings.HasPrefix(fileInfo.Name(), "."):
			result = 0
		default:
			result = fileInfo.Size()
		}
	}

	return formatSize(result, human), nil
}
