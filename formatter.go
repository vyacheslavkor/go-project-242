package code

import "fmt"

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
