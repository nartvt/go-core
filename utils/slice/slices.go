package slice

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// GetAt get element at pos of slices
func GetAt[E any](slices []E, pos int) E {
	var defaultValue E
	if len(slices) > pos {
		return slices[pos]
	}
	return defaultValue
}

// Diff Get []T in sources but not in slices
// Example:
//
//	slices.Diff(skus, products, func(p Product) string {
//		return p.TkSku
//	})
func Diff[T comparable, E any](sources []T, slices []E, f func(v E) T) []T {
	var diffStr []T
	var defaultValue T
	m := map[T]int{}

	for _, doc := range sources {
		m[doc] = 1
	}
	for _, s := range slices {
		if k := f(s); k != defaultValue {
			m[k]++
		}
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}
	return diffStr
}

// IsEmpty check slices are empty
func IsEmpty[U any](slices []U) bool {
	return len(slices) == 0
}

// IsNotEmpty check slices are not empty
func IsNotEmpty[U any](slices []U) bool {
	return len(slices) > 0
}

// Map Convert []U to []V with func transform
// Example:
//
//	slices.Map(products, func(p Product) []int64 {
//		return p.ID
//	})
func Map[U any, V any](sources []U, transfer func(u U) V) []V {
	result := make([]V, len(sources))
	for i, el := range sources {
		result[i] = transfer(el)
	}
	return result
}

// MapIndex Convert []U to []V with func transform
// Map handle index
// Example:
//
//	slices.Map2(products, func(p Product, index int) []int64 {
//		return p.ID
//	})
func MapIndex[U any, V any](sources []U, transform func(U, int) V) []V {
	result := make([]V, len(sources))
	for i, el := range sources {
		result[i] = transform(el, i)
	}
	return result
}

// MapDistinct Map and DISTINCT
func MapDistinct[U any, V comparable](sources []U, transfer func(u U) V) []V {
	result := make([]V, 0, len(sources))
	seen := make(map[V]struct{}, len(sources))
	for _, el := range sources {
		key := transfer(el)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, key)
	}
	return result
}

// FlatMap manipulates a slice and transforms and flattens it to a slice of another type.
func FlatMap[U any, V any](sources []U, transform func(U) []V) []V {
	result := make([]V, 0, len(sources))

	for _, item := range sources {
		result = append(result, transform(item)...)
	}

	return result
}

// FlatMapIndex manipulates a slice and transforms and flattens it to a slice of another type.
// FlatMap handle index
func FlatMapIndex[U any, V any](sources []U, transform func(U, int) []V) []V {
	result := make([]V, 0, len(sources))

	for i, item := range sources {
		result = append(result, transform(item, i)...)
	}

	return result
}

// Key transforms a slice or an array of structs to a map based on a pivot callback
// Example:
//
//	slices.Key(products, func(p Product) int64 {
//		return p.ID
//	})
func Key[U any, K comparable](sources []U, transform func(u U) K) map[K]U {
	m := make(map[K]U)
	for _, u := range sources {
		k := transform(u)
		m[k] = u
	}
	return m
}

// KeyBy transforms a slice or an array of structs to a map based on a pivot callback
// Example:
//
//	slices.KeyBy(products, func(p Product) (int64, string) {
//		return p.ID, p.TkSku
//	})
func KeyBy[U any, K comparable, V any](sources []U, transform func(u U) (K, V)) map[K]V {
	m := make(map[K]V)
	for _, u := range sources {
		k, v := transform(u)
		m[k] = v
	}
	return m
}

// ToFlatMap convert []U to map[K][]V
// Example:
//
//	collection.MakeAndMergeMap(products, func(p Product) (int64, string) {
//		return p.ID, p.TkSku
//	})
func ToFlatMap[U any, K comparable, V any](sources []U, transform func(u U) (K, V)) map[K][]V {
	m := make(map[K][]V)
	for i, u := range sources {
		k, v := transform(u)
		if _, ok := m[k]; !ok {
			m[k] = make([]V, 0, len(sources)-i)
		}
		m[k] = append(m[k], v)
	}
	return m
}

// Filter iterating over elements of a collection, returning an array of all elements predicate returns truthy for.
func Filter[U any](slices []U, predicate func(U) bool) []U {
	result := make([]U, 0, len(slices))
	for _, item := range slices {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

// FilterIndex Filter with predicate index
// Filter handle index
func FilterIndex[U any](slices []U, predicate func(U, int) bool) []U {
	result := make([]U, 0, len(slices))
	for i, item := range slices {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// FilterMap returns a slice which obtained after both filtering and mapping using the given callback function.
// The callback function should return two values:
//   - the result of the mapping operation and
//   - whether the result element should be included or not.
//
// Play: https://go.dev/play/p/-AuYXfy7opz
func FilterMap[T any, R any](slices []T, filterTrans func(item T) (bool, R)) []R {
	result := make([]R, 0, len(slices))
	for _, item := range slices {
		if ok, r := filterTrans(item); ok {
			result = append(result, r)
		}
	}

	return result
}

// Uniq returns a duplicate-free version of an array, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the array.
func Uniq[T comparable](slices []T) []T {
	result := make([]T, 0, len(slices))
	seen := make(map[T]struct{}, len(slices))

	for _, item := range slices {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// UniqBy returns a duplicate-free version of an array, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is
// invoked for each element in array to generate the criterion by which uniqueness is computed.
func UniqBy[T any, U comparable](slices []T, iteratee func(T) U) []T {
	result := make([]T, 0, len(slices))
	seen := make(map[U]struct{}, len(slices))

	for _, item := range slices {
		key := iteratee(item)

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, item)
	}

	return result
}

// GroupBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
func GroupBy[U any, K comparable](slices []U, iteratee func(U) K) map[K][]U {
	result := make(map[K][]U)

	for _, item := range slices {
		key := iteratee(item)
		result[key] = append(result[key], item)
	}

	return result
}

// ForEach iterates over elements of collection and invokes iteratee for each element.
func ForEach[T any](collection []T, iteratee func(T)) {
	for _, item := range collection {
		iteratee(item)
	}
}

// ForEachIndex iterates over elements of collection and invokes iteratee for each element.
// ForEach handle index
func ForEachIndex[T any](collection []T, iteratee func(T, int)) {
	for i, item := range collection {
		iteratee(item, i)
	}
}

// Reduce reduces collection to a value which is the accumulated result of running each element in collection
// through accumulator, where each successive invocation is supplied the return value of the previous.
func Reduce[T any, R any](collection []T, accumulator func(R, T, int) R, initial R) R {
	for i, item := range collection {
		initial = accumulator(initial, item, i)
	}

	return initial
}

// GroupMapBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
func GroupMapBy[U any, K comparable, V any](slices []U, trans func(U) (K, V)) map[K][]V {
	result := make(map[K][]V)

	for _, item := range slices {
		key, value := trans(item)
		result[key] = append(result[key], value)
	}

	return result
}

// GroupFlatMapBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
func GroupFlatMapBy[U any, K comparable, V any](slices []U, trans func(U) (K, []V)) map[K][]V {
	result := make(map[K][]V)

	for _, item := range slices {
		key, values := trans(item)
		result[key] = append(result[key], values...)
	}

	return result
}

// Contains Check item in slice T type
func Contains[T comparable](slice []T, item T) bool {
	for _, item2 := range slice {
		if item2 == item {
			return true
		}
	}
	return false
}

// ContainsBy returns true if predicate function return true.
func ContainsBy[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}

	return false
}

// First iterates over elements of collection, returning first element returns truthy for.
func First[U any](slices []U, predicate func(U) bool) U {
	var defaultValue U
	for _, item := range slices {
		if predicate(item) {
			return item
		}
	}
	return defaultValue
}

// Every returns true if all elements of a subset are contained into a collection or if the subset is empty.
func Every[T comparable](collection []T, subset []T) bool {
	for _, elem := range subset {
		if !Contains(collection, elem) {
			return false
		}
	}

	return true
}

// EveryBy returns true if the predicate returns true for all of the elements in the collection or if the collection is empty.
func EveryBy[T any](collection []T, predicate func(item T) bool) bool {
	for _, v := range collection {
		if !predicate(v) {
			return false
		}
	}

	return true
}

// Flatten returns an array a single level deep.
// Play: https://go.dev/play/p/rbp9ORaMpjw
func Flatten[T any](collection [][]T) []T {
	totalLen := 0
	for i := range collection {
		totalLen += len(collection[i])
	}

	result := make([]T, 0, totalLen)
	for i := range collection {
		result = append(result, collection[i]...)
	}

	return result
}

// Min search the minimum value of a collection.
// Returns zero value when collection is empty.
func Min[T constraints.Ordered](collection []T) T {
	var min T

	if len(collection) == 0 {
		return min
	}

	min = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item < min {
			min = item
		}
	}

	return min
}

// MinBy search the minimum value of a collection using the given comparison function.
// If several values of the collection are equal to the smallest value, returns the first such value.
// Returns zero value when collection is empty.
func MinBy[T any](collection []T, comparison func(a T, b T) bool) T {
	var min T

	if len(collection) == 0 {
		return min
	}

	min = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if comparison(item, min) {
			min = item
		}
	}

	return min
}

// Max searches the maximum value of a collection.
// Returns zero value when collection is empty.
func Max[T constraints.Ordered](collection []T) T {
	var max T

	if len(collection) == 0 {
		return max
	}

	max = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item > max {
			max = item
		}
	}

	return max
}

// MaxBy search the maximum value of a collection using the given comparison function.
// If several values of the collection are equal to the greatest value, returns the first such value.
// Returns zero value when collection is empty.
func MaxBy[T any](collection []T, comparison func(a T, b T) bool) T {
	var max T

	if len(collection) == 0 {
		return max
	}

	max = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if comparison(item, max) {
			max = item
		}
	}

	return max
}

// Nth returns the element at index `nth` of collection. If `nth` is negative, the nth element
// from the end is returned. An error is returned when nth is out of slice bounds.
func Nth[T any, N constraints.Integer](collection []T, nth N) (T, error) {
	n := int(nth)
	l := len(collection)
	if n >= l || -n > l {
		var t T
		return t, fmt.Errorf("nth: %d out of slice bounds", n)
	}

	if n >= 0 {
		return collection[n], nil
	}
	return collection[l+n], nil
}

func Make[T any](values ...T) []T {
	res := make([]T, 0, len(values))
	if len(values) > 0 {
		res = append(res, values...)
	}
	return res
}
