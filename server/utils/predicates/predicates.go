package predicates

// One returns true if the conditional function returns true for any element in the array
func One[T any](slice []T, cond func(a T) bool) bool {
	for _, v := range slice {
		if cond(v) {
			return true
		}
	}
	return false
}
