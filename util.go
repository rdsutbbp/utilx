package utilx

// StringInSlice return true if string in the slice `list`
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}


