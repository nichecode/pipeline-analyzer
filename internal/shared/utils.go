package shared

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// GenerateID creates a unique identifier for analysis results
func GenerateID(prefix string) string {
	timestamp := time.Now().Unix()
	hash := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", prefix, timestamp))))
	return fmt.Sprintf("%s-%s", prefix, hash[:8])
}

// NormalizeFileName creates a safe filename from any string
func NormalizeFileName(name string) string {
	// Replace problematic characters
	replacements := map[string]string{
		"/":  "-",
		"\\": "-",
		":":  "-",
		"*":  "-",
		"?":  "-",
		"\"": "-",
		"<":  "-",
		">":  "-",
		"|":  "-",
		" ":  "-",
		"_":  "-",
	}
	
	normalized := name
	for old, new := range replacements {
		normalized = strings.ReplaceAll(normalized, old, new)
	}
	
	// Remove multiple consecutive dashes
	for strings.Contains(normalized, "--") {
		normalized = strings.ReplaceAll(normalized, "--", "-")
	}
	
	// Trim dashes from start and end
	normalized = strings.Trim(normalized, "-")
	
	// Ensure lowercase
	normalized = strings.ToLower(normalized)
	
	// Ensure it's not empty
	if normalized == "" {
		normalized = "unnamed"
	}
	
	return normalized
}

// GetRelativePathSafe safely gets a relative path
func GetRelativePathSafe(basepath, targetpath string) string {
	rel, err := filepath.Rel(basepath, targetpath)
	if err != nil {
		return targetpath
	}
	return rel
}

// TruncateString safely truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	
	if maxLen <= 3 {
		return s[:maxLen]
	}
	
	return s[:maxLen-3] + "..."
}

// SanitizeForMarkdown escapes special markdown characters
func SanitizeForMarkdown(s string) string {
	replacements := map[string]string{
		"*":  "\\*",
		"_":  "\\_",
		"#":  "\\#",
		"`":  "\\`",
		"[":  "\\[",
		"]":  "\\]",
		"(":  "\\(",
		")":  "\\)",
		"!":  "\\!",
		"|":  "\\|",
		"\\": "\\\\",
	}
	
	result := s
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	
	return result
}

// FormatDuration formats a duration for human readability
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// CalculateComplexityScore calculates a complexity score for analysis items
func CalculateComplexityScore(factors map[string]int, weights map[string]float64) float64 {
	if len(factors) == 0 {
		return 0.0
	}
	
	totalScore := 0.0
	totalWeight := 0.0
	
	for factor, count := range factors {
		weight := weights[factor]
		if weight == 0 {
			weight = 1.0 // Default weight
		}
		totalScore += float64(count) * weight
		totalWeight += weight
	}
	
	if totalWeight == 0 {
		return 0.0
	}
	
	return totalScore / totalWeight
}

// GroupItems groups items by a key function
func GroupItems[T any, K comparable](items []T, keyFunc func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, item := range items {
		key := keyFunc(item)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// FilterItems filters items based on a predicate function
func FilterItems[T any](items []T, predicate func(T) bool) []T {
	var filtered []T
	for _, item := range items {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// MapItems transforms items using a mapping function
func MapItems[T any, R any](items []T, mapFunc func(T) R) []R {
	var mapped []R
	for _, item := range items {
		mapped = append(mapped, mapFunc(item))
	}
	return mapped
}

// UniqueStrings returns unique strings from a slice
func UniqueStrings(items []string) []string {
	seen := make(map[string]bool)
	var unique []string
	
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			unique = append(unique, item)
		}
	}
	
	return unique
}

// ContainsString checks if a string slice contains a specific string
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// MergeStringSlices merges multiple string slices, removing duplicates
func MergeStringSlices(slices ...[]string) []string {
	seen := make(map[string]bool)
	var result []string
	
	for _, slice := range slices {
		for _, item := range slice {
			if !seen[item] {
				seen[item] = true
				result = append(result, item)
			}
		}
	}
	
	return result
}

// SplitCommand safely splits a command string into parts
func SplitCommand(command string) []string {
	// Simple split that handles quoted strings
	var parts []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune
	
	for _, r := range command {
		switch r {
		case '"', '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = r
			} else if r == quoteChar {
				inQuotes = false
			}
			current.WriteRune(r)
		case ' ', '\t':
			if inQuotes {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	
	return parts
}

// ExtractCommandName extracts the command name from a command string
func ExtractCommandName(command string) string {
	parts := SplitCommand(strings.TrimSpace(command))
	if len(parts) == 0 {
		return ""
	}
	
	// Handle common prefixes
	cmdName := parts[0]
	
	// Remove path if present
	if strings.Contains(cmdName, "/") {
		cmdName = filepath.Base(cmdName)
	}
	
	// Remove common prefixes
	prefixes := []string{"sudo", "time", "timeout", "nice"}
	for _, prefix := range prefixes {
		if cmdName == prefix && len(parts) > 1 {
			cmdName = parts[1]
			if strings.Contains(cmdName, "/") {
				cmdName = filepath.Base(cmdName)
			}
			break
		}
	}
	
	return cmdName
}

// IsValidVersion checks if a version string is valid
func IsValidVersion(version string) bool {
	validVersions := []string{"1", "2", "2.1", "2.2", "2.6", "3"}
	for _, v := range validVersions {
		if version == v {
			return true
		}
	}
	return false
}

// CompareVersions compares two version strings
func CompareVersions(v1, v2 string) int {
	// Simple version comparison for task file versions
	versions := map[string]int{
		"1":   1,
		"2":   2,
		"2.1": 21,
		"2.2": 22,
		"2.6": 26,
		"3":   3,
	}
	
	val1, ok1 := versions[v1]
	val2, ok2 := versions[v2]
	
	if !ok1 || !ok2 {
		return strings.Compare(v1, v2)
	}
	
	if val1 < val2 {
		return -1
	} else if val1 > val2 {
		return 1
	}
	return 0
}

// GenerateProgressBar creates a simple text progress bar
func GenerateProgressBar(current, total int, width int) string {
	if total == 0 {
		return strings.Repeat("█", width)
	}
	
	progress := float64(current) / float64(total)
	filled := int(progress * float64(width))
	
	if filled > width {
		filled = width
	}
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("%s %d/%d (%.1f%%)", bar, current, total, progress*100)
}

// CalculateLevenshteinDistance calculates the edit distance between two strings
func CalculateLevenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	
	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
	}
	
	// Initialize first row and column
	for i := 0; i <= len(a); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}
	
	// Fill the matrix
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			
			matrix[i][j] = min3(
				matrix[i-1][j]+1,     // deletion
				matrix[i][j-1]+1,     // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}
	
	return matrix[len(a)][len(b)]
}

// min3 returns the minimum of three integers
func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}