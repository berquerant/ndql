package regexpx

import "strings"

// Convert SQL 'Like' expression into regexp.
func LikeToRegexpString(s string, escape rune) string {
	var (
		b       strings.Builder
		escaped bool
	)
	for _, c := range s {
		switch {
		case escaped:
			escaped = false
			b.WriteRune(c)
		case c == escape:
			escaped = true
		case c == '%':
			b.WriteString(".*")
		case c == '_':
			b.WriteRune('.')
		default:
			b.WriteRune(c)
		}
	}
	return b.String()
}

func LikeToRegexpStringDefault(s string) string {
	return LikeToRegexpString(s, '|')
}
