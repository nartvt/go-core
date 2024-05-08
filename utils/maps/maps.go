package maps

import (
	"github.com/nartvt/go-core/utils/slice"
	"github.com/nartvt/go-core/utils/typez"
)

// IsEmpty check map is empty
func IsEmpty[K comparable, V any](in map[K]V) bool {
	return len(in) == 0
}

// IsNotEmpty check map is not empty
func IsNotEmpty[K comparable, V any](in map[K]V) bool {
	return len(in) > 0
}

// Keys create an array of the map keys.
// Play: https://go.dev/play/p/Uu11fHASqrU
func Keys[K comparable, V any](in map[K]V) []K {
	result := make([]K, 0, len(in))

	for k := range in {
		result = append(result, k)
	}

	return result
}

// Values create an array of the map values.
// Play: https://go.dev/play/p/nnRTQkzQfF6
func Values[K comparable, V any](in map[K]V) []V {
	result := make([]V, 0, len(in))

	for _, v := range in {
		result = append(result, v)
	}

	return result
}

// FlatValues creates an array of the map values.
// Play: https://go.dev/play/p/nnRTQkzQfF6
func FlatValues[K comparable, V any](in map[K][]V) []V {
	result := make([]V, 0, len(in))

	for _, v := range in {
		result = append(result, v...)
	}

	return result
}

// Entries transform a map into an array of key/value pairs.
// Play:
func Entries[K comparable, V any](in map[K]V) []typez.Entry[K, V] {
	entries := make([]typez.Entry[K, V], 0, len(in))

	for k, v := range in {
		entries = append(entries, typez.Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}

	return entries
}

// ToPairs transforms a map into an array of key/value pairs.
// Alias of Entries().
// Play: https://go.dev/play/p/3Dhgx46gawJ
func ToPairs[K comparable, V any](in map[K]V) []typez.Entry[K, V] {
	return Entries(in)
}

// FromEntries transforms an array of key/value pairs into a map.
// Play: https://go.dev/play/p/oIr5KHFGCEN
func FromEntries[K comparable, V any](entries []typez.Entry[K, V]) map[K]V {
	out := map[K]V{}

	for _, v := range entries {
		out[v.Key] = v.Value
	}

	return out
}

// FromPairs transforms an array of key/value pairs into a map.
// Alias of FromEntries().
// Play: https://go.dev/play/p/oIr5KHFGCEN
func FromPairs[K comparable, V any](entries []typez.Entry[K, V]) map[K]V {
	return FromEntries(entries)
}

// ToSlice transforms a map into a slice based on specific iteratee
// Play: https://go.dev/play/p/ZuiCZpDt6LD
func ToSlice[K comparable, V any, R any](in map[K]V, iteratee func(key K, value V) R) []R {
	result := make([]R, 0, len(in))

	for k, v := range in {
		result = append(result, iteratee(k, v))
	}

	return result
}

// PickBy returns the same map type filtered by given predicate.
// Play: https://go.dev/play/p/kdg8GR_QMmf
func PickBy[K comparable, V any](in map[K]V, predicate func(key K, value V) bool) map[K]V {
	r := map[K]V{}
	for k, v := range in {
		if predicate(k, v) {
			r[k] = v
		}
	}
	return r
}

// PickByKeys returns a same map type filtered by given keys.
// Play: https://go.dev/play/p/R1imbuci9qU
func PickByKeys[K comparable, V any](in map[K]V, keys []K) map[K]V {
	r := map[K]V{}
	for k, v := range in {
		if slice.Contains(keys, k) {
			r[k] = v
		}
	}
	return r
}

// PickByValues returns same map type filtered by given values.
// Play: https://go.dev/play/p/1zdzSvbfsJc
func PickByValues[K comparable, V comparable](in map[K]V, values []V) map[K]V {
	r := map[K]V{}
	for k, v := range in {
		if slice.Contains(values, v) {
			r[k] = v
		}
	}
	return r
}

func GetValue[K comparable, V comparable](in map[K]V, key K) V {
	var defaultValue V
	if v, ok := in[key]; ok {
		return v
	}
	return defaultValue
}

func GetCastValue[K comparable, V any](in map[K]any, key K) V {
	var defaultValue V
	if v, ok := in[key]; ok {
		return typez.Cast[V](v)
	}
	return defaultValue
}
