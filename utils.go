package inversify

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getFirstStringArgumentOrEmpty(names []string) string {
	if len(names) > 0 {
		return names[0]
	}
	return ""
}
