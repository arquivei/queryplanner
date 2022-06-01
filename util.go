package queryplanner

func isInArray(s string, arr []string) bool {
	for _, item := range arr {
		if s == item {
			return true
		}
	}
	return false
}
