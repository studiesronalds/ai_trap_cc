package output

import (
	"fmt"
	"io"
	"math"
	"strings"

	"aisim/internal/model"
)

// Heatmap2D prints a 2D ASCII heatmap of the wedge over a grid.
// Rows = p1 (y-axis), Cols = p2 (x-axis). Darker = larger wedge.
func Heatmap2D(w io.Writer, grid [][]model.Sweep2DPoint) {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return
	}
	p1Name := grid[0][0].P1Name
	p2Name := grid[0][0].P2Name

	// Compute range
	minW, maxW := math.Inf(1), math.Inf(-1)
	for _, row := range grid {
		for _, cell := range row {
			if cell.Wedge < minW { minW = cell.Wedge }
			if cell.Wedge > maxW { maxW = cell.Wedge }
		}
	}

	shades := []string{" ", "░", "▒", "▓", "█"}

	cell := func(wedge float64) string {
		if maxW <= minW {
			return shades[0]
		}
		idx := int((wedge - minW) / (maxW - minW) * float64(len(shades)-1))
		if idx < 0 { idx = 0 }
		if idx >= len(shades) { idx = len(shades) - 1 }
		return shades[idx]
	}

	// Header: p2 axis values
	fmt.Fprintf(w, "\nWedge heatmap: rows=%s, cols=%s  (darker = larger wedge)\n", p1Name, p2Name)
	fmt.Fprintf(w, "Range: %.4f – %.4f\n\n", minW, maxW)

	// Column header
	nCols := len(grid[0])
	colWidth := 6
	fmt.Fprintf(w, "%-8s", p1Name+`\`+p2Name)
	for j := 0; j < nCols; j++ {
		fmt.Fprintf(w, " %*.*f", colWidth, 2, grid[0][j].P2)
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s\n", strings.Repeat("─", 8+nCols*(colWidth+1)))

	for _, row := range grid {
		fmt.Fprintf(w, "%-8.3f", row[0].P1)
		for _, c := range row {
			// print shade + value
			fmt.Fprintf(w, " %s%-*s", cell(c.Wedge), colWidth-1, fmt.Sprintf("%.3f", c.Wedge))
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
}
