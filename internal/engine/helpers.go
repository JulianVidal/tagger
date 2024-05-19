package engine

import "errors"

func delItemsFromSlice[S ~[]I, I comparable](s S, items ...I) (S, error) {
	deleted := 0

	for _, item := range items {
		for i, value := range s {
			if item == value {
				s[i] = s[len(s)-1]
				s = s[:len(s)-1]
				deleted += 1
				break
			}
		}
	}
	if len(items) == deleted {
		return s, nil
	}

	return nil, errors.New("Not all items were deleted")
}

// TODO: How should it deal if the same key has two different values, which should it pick?
func mapUnion[M map[K]V, K comparable, V comparable](a M, b M) {
	for k, v := range b {
		if _, exist := a[k]; exist {
			if a[k] != v {
				panic("Map have the same key but different value.")
			}
		}
		a[k] = v
	}
}

func mapKeys[M map[K]V, K comparable, V any](m M) []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func mapValues[M map[K]V, K comparable, V any](m M) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}
