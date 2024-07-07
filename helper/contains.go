package helper

func Contains[K comparable](slice []K, str K) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
