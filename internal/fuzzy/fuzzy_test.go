package fuzzy

import "testing"

func TestMatchScore(t *testing.T) {
	tests := []struct {
		name   string
		target string
		query  string
		wantOK bool
	}{
		{"empty query matches everything", "anything/at/all.go", "", true},
		{"empty target with empty query", "", "", true},
		{"exact match", "save", "save", true},
		{"case-insensitive", "Save File", "save", true},
		{"single char", "editor.go", "e", true},
		{"in-order subsequence", "src/components/editor.go", "sced", true},
		{"query longer than target", "edit", "editor", false},
		{"out of order fails", "src/components/editor.go", "eds", false},
		{"missing char fails", "editor.go", "editors", false},
		{"no overlap fails", "editor.go", "xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, _ := MatchScore(tt.target, tt.query)
			if ok != tt.wantOK {
				t.Fatalf("MatchScore(%q, %q) ok = %v, want %v", tt.target, tt.query, ok, tt.wantOK)
			}
		})
	}
}

// TestMatchScoreOrdering verifies the ranking incentives that make fuzzy
// search feel good: consecutive matches beat scattered ones, matches after a
// path separator beat mid-word ones, and camelCase boundaries beat mid-word
// matches.
func TestMatchScoreOrdering(t *testing.T) {
	score := func(target, query string) int {
		ok, s := MatchScore(target, query)
		if !ok {
			t.Fatalf("MatchScore(%q, %q) did not match", target, query)
		}
		return s
	}

	if got, want := score("editor.go", "ed"), score("editor.go", "eo"); got <= want {
		t.Errorf("consecutive match (%d) should outscore scattered (%d)", got, want)
	}

	if got, want := score("a/bcd", "bc"), score("a/bcd", "cd"); got <= want {
		t.Errorf("match after path separator (%d) should outscore mid-word (%d)", got, want)
	}

	if got, want := score("SaveFile", "sf"), score("SaveFile", "se"); got <= want {
		t.Errorf("camelCase boundary match (%d) should outscore mid-word (%d)", got, want)
	}
}
