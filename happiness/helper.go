package happiness

func SliceContains[T comparable](haystack []T, needle T) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}
