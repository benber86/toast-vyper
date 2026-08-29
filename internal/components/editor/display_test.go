package editor

import (
	"testing"
)

// Regression tests for the grapheme-cluster width model: emoji with a
// variation selector (e.g. "⚠️" = U+26A0 + U+FE0F) render as ONE cluster two
// cells wide, but per-rune width counting used to treat them as 1 cell, making
// every following position (cursor X, mouse clicks, viewport, wrap breaks) off
// by one. This was the "caution emoji in the README shifts the cursor" bug.

func TestDisplayColumnAtByte_VS16EmojiCluster(t *testing.T) {
	line := "> ⚠️ This"
	afterEmoji := len("> ⚠️")
	// '>' 1 + ' ' 1 + "⚠️" 2 = 4 cells up to the space after the cluster.
	if got := displayColumnAtByte(line, afterEmoji, 4); got != 4 {
		t.Fatalf("displayColumnAtByte after ⚠️ = %d, want 4 (two-cell cluster)", got)
	}
	// Whole line: 4 + ' ' 1 + "This" 4 = 9.
	if got := displayColumnAtByte(line, len(line), 4); got != 9 {
		t.Fatalf("line width = %d, want 9", got)
	}
}

func TestDisplayWidthForByteRange_VS16EmojiCluster(t *testing.T) {
	line := "⚠️abc"
	// Width of the cluster alone is 2 cells, not 1.
	if got := displayWidthForByteRange(line, 0, len("⚠️"), 4); got != 2 {
		t.Fatalf("width of ⚠️ cluster = %d, want 2", got)
	}
	if got := displayWidthForByteRange(line, 0, len(line), 4); got != 5 {
		t.Fatalf("width of %q = %d, want 5", line, got)
	}
}

func TestByteColForDisplayOffset_VS16EmojiCluster(t *testing.T) {
	line := "> ⚠️ This"
	afterEmoji := len("> ⚠️")
	// Cell 3 is the second half of the 2-cell cluster → maps before the cluster.
	if c := byteColForDisplayOffset(line, 0, 3, 4); c != len("> ") {
		t.Fatalf("byteColForDisplayOffset(3) = %d, want %d (inside cluster → before it)", c, len("> "))
	}
	// Cell 4 (right after the cluster) maps to the byte after it.
	if c := byteColForDisplayOffset(line, 0, 4, 4); c != afterEmoji {
		t.Fatalf("byteColForDisplayOffset(4) = %d, want %d (right after cluster)", c, afterEmoji)
	}
}

func TestWordWrapChunks_VS16EmojiCluster(t *testing.T) {
	// "⚠️a" fills the 3-cell chunk exactly (cluster 2 cells + 'a' 1 cell).
	// The old per-rune width counted the whole 5-cell line as fitting in 3.
	got := wordWrapChunks("⚠️ab", 3)
	want := []int{0, len("⚠️a")}
	if len(got) != len(want) {
		t.Fatalf("chunks = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("chunks = %v, want %v", got, want)
		}
	}
}

func TestCursorScreenPositionAfterVS16Emoji(t *testing.T) {
	m := newThemedTestModel(t, "⚠️abc\n")
	m.cursor = cursorPos{line: 0, col: len("⚠️")} // right after the cluster

	v := m.View()
	if v.Cursor == nil {
		t.Fatal("expected cursor")
	}
	// "⚠️" occupies two cells, so the cursor lands on the third content cell
	// (over 'a'), matching where the next typed character will appear.
	wantX := m.gutterWidth + 2
	if v.Cursor.Position.X != wantX {
		t.Fatalf("cursor X after ⚠️ = %d, want %d", v.Cursor.Position.X, wantX)
	}
}
