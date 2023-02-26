package predicates

func One[T any](slice []T, cond func(a T) bool) bool {
	for _, v := range slice {
		if cond(v) {
			return true
		}
	}
	return false
}
