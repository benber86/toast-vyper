// Package fuzzy provides subsequence-based fuzzy matching with a
// position-and-structure-aware scoring function. It is shared by overlay
// components that search lists of strings (quick-open files, command palette
// commands, ...).
package fuzzy

import "strings"

// MatchScore returns true with a score if query matches target (case-insensitive).
// Characters in query must appear in order within target. The score favours
// consecutive matches, matches after separators, and matches at word boundaries.
func MatchScore(target, query string) (bool, int) {
	if query == "" {
		return true, 0
	}

	targetLower := strings.ToLower(target)
	queryLower := strings.ToLower(query)

	qi := 0
	score := 0
	lastMatch := -2 // ensure first match doesn't get a consecutive bonus

	for si := 0; si < len(targetLower) && qi < len(queryLower); si++ {
		if targetLower[si] == queryLower[qi] {
			// Consecutive match bonus
			if si == lastMatch+1 {
				score += 8
			}
			// Match after path separator bonus
			if si == 0 || targetLower[si-1] == '/' || targetLower[si-1] == '\\' {
				score += 12
			}
			// Word boundary bonus (separator chars)
			if si > 0 && (targetLower[si-1] == '_' || targetLower[si-1] == '-' ||
				targetLower[si-1] == '.' || targetLower[si-1] == ' ' ||
				targetLower[si-1] == '(' || targetLower[si-1] == '[') {
				score += 8
			}
			// Uppercase in camelCase bonus (only applies when query is lowercase
			// but target has uppercase at this position)
			if target[si] >= 'A' && target[si] <= 'Z' && queryLower[qi] == query[qi] {
				score += 6
			}
			score += 1
			qi++
			lastMatch = si
		}
	}

	if qi == len(queryLower) {
		return true, score
	}
	return false, 0
}
