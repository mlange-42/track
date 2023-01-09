package util

// Reverse reverts a slice in-place
func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Unique returns a slice with only unique elements of the original slice
func Unique[T comparable](slice []T) []T {
	mapping := map[T]bool{}
	unique := []T{}
	for _, t := range slice {
		if _, ok := mapping[t]; !ok {
			mapping[t] = true
			unique = append(unique, t)
		}
	}
	return unique
}
