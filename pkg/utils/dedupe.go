package utils

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/agnivade/levenshtein"
)

// computeSimilarity returns the normalized similarity score [0,1]
// uses Levenshtein distance to compute similarity: https://en.wikipedia.org/wiki/Levenshtein_distance
// a score of 1.0 means the strings are identical, and 0.0 means they are completely different
// Levenshtein distance is a measure of the difference between two sequences
// it is defined as the minimum number of single-character edits (insertions, deletions, or substitutions) required to change one string into the other
func computeSimilarity(a, b string) float64 {
	dist := levenshtein.ComputeDistance(a, b)
	maxLen := max(len(a), len(b))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(dist)/float64(maxLen)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// DeduplicateLogsWithWindow removes similar lines based on relative similarity threshold
// within a sliding window of the last `windowSize` lines.
func DeduplicateLogsWithWindow(input string, threshold float64, windowSize int) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var buffer bytes.Buffer
	var recentLines []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				break
			}
			if line == "" {
				break
			}
			// process the last line
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		isDuplicate := false
		for _, prev := range recentLines {
			if computeSimilarity(line, prev) >= threshold {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			buffer.WriteString(line + "\n")
			recentLines = append(recentLines, line)
			if len(recentLines) > windowSize {
				recentLines = recentLines[1:]
			}
		}

		if err != nil {
			break
		}
	}
	return buffer.String()
}
