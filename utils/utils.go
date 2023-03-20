package utils

import "os"

func CreateFile() (*os.File, error) {
	return os.Create("../Files/File1.pdf")
}

func CalculateColumnWidths(ratios []float64, totalWidth, minWidth, margin float64) []float64 {
	numColumns := len(ratios)
	remainingWidth := totalWidth - 2*margin

	// Calculate the total minimum width required for all columns
	totalMinWidth := minWidth * float64(numColumns)

	// If the total minimum width is greater than the total width, adjust the minimum width for each column accordingly
	if totalMinWidth > totalWidth {
		minWidth = totalWidth / float64(numColumns)
	}

	// Calculate the actual width of each column based on its ratio and the total width
	widths := make([]float64, numColumns)
	for i, ratio := range ratios {
		width := ratio * totalWidth
		if width < minWidth {
			width = minWidth
		}
		if width > remainingWidth {
			width = remainingWidth
		}
		remainingWidth -= width
		widths[i] = width
	}

	return widths
}
