package editor

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

const defaultTabWidth = 4

func normalizedTabWidth(width int) int {
	if width < 1 {
		return defaultTabWidth
	}
	return width
}

// firstGraphemeClusterWidth returns the first grapheme cluster of s and its
// display width in cells. Clusters are measured the way a terminal renders
// them: an emoji with a variation selector or a ZWJ sequence (e.g. "⚠️") is a
// single cluster occupying two cells. Per-rune width would count "⚠" plus the
// zero-width variation selector as one cell, leaving every following position
// (cursor, mouse clicks, viewport) off by one.
func firstGraphemeClusterWidth(s string) (string, int) {
	if s == "" {
		return "", 0
	}
	return ansi.FirstGraphemeCluster(s, ansi.GraphemeWidth)
}

// nextDisplayColumn returns the display column after a grapheme cluster,
// expanding tabs to the next tab stop.
func nextDisplayColumn(column int, cluster string, tabWidth int) int {
	if cluster == "\t" {
		tabWidth = normalizedTabWidth(tabWidth)
		return column + tabWidth - column%tabWidth
	}
	return column + ansi.StringWidth(cluster)
}

func displayColumnAtByte(text string, byteCol, tabWidth int) int {
	byteCol = clampByteCol(text, byteCol)
	column := 0
	for offset := 0; offset < byteCol; {
		cluster, _ := firstGraphemeClusterWidth(text[offset:])
		next := offset + len(cluster)
		if next > byteCol {
			// byteCol cuts inside a grapheme cluster (e.g. between a base rune
			// and its variation selector). The cluster is one atomic glyph, so
			// its width applies only at or after `next`.
			break
		}
		column = nextDisplayColumn(column, cluster, tabWidth)
		offset = next
	}
	return column
}

func displayWidthForByteRange(text string, start, end, tabWidth int) int {
	start = clampByteCol(text, start)
	end = clampByteCol(text, end)
	if end < start {
		end = start
	}
	return displayColumnAtByte(text, end, tabWidth) - displayColumnAtByte(text, start, tabWidth)
}

// byteColForDisplayOffset maps an offset from start's display column back to
// a byte boundary. Cells occupied by a wide cluster or expanded tab map to the
// position before that character; its right boundary maps after it.
func byteColForDisplayOffset(text string, start, displayOffset, tabWidth int) int {
	start = clampByteCol(text, start)
	if displayOffset <= 0 {
		return start
	}
	target := displayColumnAtByte(text, start, tabWidth) + displayOffset
	column := displayColumnAtByte(text, start, tabWidth)
	previous := start
	for offset := start; offset < len(text); {
		cluster, _ := firstGraphemeClusterWidth(text[offset:])
		next := offset + len(cluster)
		nextColumn := nextDisplayColumn(column, cluster, tabWidth)
		if nextColumn > target {
			return previous
		}
		if nextColumn == target {
			return next
		}
		previous = next
		offset = next
		column = nextColumn
	}
	return len(text)
}

// byteColAtOrAfterDisplayColumn finds a safe horizontal viewport boundary.
// When the requested column falls inside a tab or wide cluster, it advances
// past that character so the cursor is guaranteed to become visible.
func byteColAtOrAfterDisplayColumn(text string, target, tabWidth int) int {
	if target <= 0 {
		return 0
	}
	column := 0
	for offset := 0; offset < len(text); {
		cluster, _ := firstGraphemeClusterWidth(text[offset:])
		next := offset + len(cluster)
		nextColumn := nextDisplayColumn(column, cluster, tabWidth)
		if column >= target {
			return offset
		}
		if nextColumn >= target {
			return next
		}
		offset = next
		column = nextColumn
	}
	return len(text)
}

func expandTabs(text string, startColumn, tabWidth int) string {
	if !strings.ContainsRune(text, '\t') {
		return text
	}
	var out strings.Builder
	column := startColumn
	for offset := 0; offset < len(text); {
		cluster, _ := firstGraphemeClusterWidth(text[offset:])
		if cluster == "\t" {
			next := nextDisplayColumn(column, cluster, tabWidth)
			out.WriteString(strings.Repeat(" ", next-column))
			column = next
		} else {
			out.WriteString(cluster)
			column = nextDisplayColumn(column, cluster, tabWidth)
		}
		offset += len(cluster)
	}
	return out.String()
}
