package lib

func Unique(a []string) bool {
	seen := make(map[string]struct{})
	for _, i := range a {
		if _, ok := seen[i]; ok {
			return false
		}
		seen[i] = struct{}{}
	}
	return true
}
