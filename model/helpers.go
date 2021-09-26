package model

// StringPtr obtains a pointer of provided string input
func StringPtr(s string) *string {
	return &s
}

// StringSliceContains determines if str exists in slice
func StringSliceContains(slice *[]string, str string) bool {
	if slice == nil {
		return false
	}

	for _, s := range *slice {
		if s == str {
			return true
		}
	}
	return false
}
