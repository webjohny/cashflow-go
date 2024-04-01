package helper

func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func ContainsInt(slice []int, str int) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
