package formatter

import "fmt"

var units = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// FormatSize returns a string representation of size in bytes or human-readable form.
func FormatSize(size int64, human bool) string {
	if human {
		return formatToHuman(size)
	}
	return fmt.Sprintf("%d%s", size, units[0])
}

// FormatOutput returns a string representation of size and path.
func FormatOutput(sizeStr, path string) string {
	return fmt.Sprintf("%s\t%s", sizeStr, path)
}

func formatToHuman(size int64) string {
	sizeFloat := float64(size)
	const sizeStepMultiplier float64 = 1024

	for _, unit := range units[:len(units)-1] {
		if sizeFloat < sizeStepMultiplier {
			if unit == "B" {
				return fmt.Sprintf("%.0f%s", sizeFloat, unit)
			}

			return fmt.Sprintf("%.1f%s", sizeFloat, unit)
		}

		sizeFloat /= sizeStepMultiplier
	}

	return fmt.Sprintf("%.1f%s", sizeFloat, units[len(units)-1])
}
