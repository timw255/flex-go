package util

// Index ...
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Contains ...
func Contains(a []string, x string) bool {
	return Index(a, x) >= 0
}

// Any ...
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// All ...
func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// Filter ...
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Map ...
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// Delete ...
func Delete(a []string, x string) []string {
	i := Index(a, x)

	if i >= 0 {
		copy(a[i:], a[i+1:])
		a[len(a)-1] = "" // or the zero value of T
		a = a[:len(a)-1]
	}

	return a
}
