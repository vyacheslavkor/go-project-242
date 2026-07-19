package sizefmt

import "fmt"

var units = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// Format returns a string representation of size in bytes or human-readable form.
func Format(size int64, human bool) string {
	if human {
		return toHuman(size)
	}
	return fmt.Sprintf("%d%s", size, units[0])
}

// PrepareOutput returns a string representation of size and path.
func PrepareOutput(sizeStr, path string) string {
	return fmt.Sprintf("%s\t%s", sizeStr, path)
}

func toHuman(size int64) string {
	humanSize := float64(size)
	const unitStep float64 = 1024

	for i, unit := range units[:len(units)-1] {
		if humanSize < unitStep {
			if i == 0 {
				return fmt.Sprintf("%.0f%s", humanSize, unit)
			}

			return fmt.Sprintf("%.1f%s", humanSize, unit)
		}

		humanSize /= unitStep
	}

	return fmt.Sprintf("%.1f%s", humanSize, units[len(units)-1])
}
