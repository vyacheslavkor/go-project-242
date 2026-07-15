package sizefmt

import "fmt"

var units = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// WithUnits returns a string representation of size in bytes or human-readable form.
func WithUnits(size int64, human bool) string {
	if human {
		return toHuman(size)
	}
	return fmt.Sprintf("%d%s", size, units[0])
}

// ToOutput returns a string representation of size and path.
func ToOutput(sizeStr, path string) string {
	return fmt.Sprintf("%s\t%s", sizeStr, path)
}

func toHuman(size int64) string {
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
