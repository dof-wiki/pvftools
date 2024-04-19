package utils

func RemoveWrapper(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == '[' && s[len(s)-1] == ']' {
		return s[1 : len(s)-1]
	}
	return s
}
