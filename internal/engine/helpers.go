package engine

import "errors"

func delItemFromSlice[S ~[]I, I comparable](s S, a I) (S, error) {
	for i, b := range s {
		if a == b {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
			return s, nil
		}
	}
	return nil, errors.New("Item to delete not found in slice")
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
