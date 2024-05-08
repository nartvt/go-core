package typez

import (
	"reflect"

	"github.com/nartvt/go-core/utils/slice"
)

// Entry defines a key/value pairs.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Tuple2 is a group of 2 elements (pair).
type Tuple2[A any, B any] struct {
	A A
	B B
}

// Unpack returns values contained in tuple.
func (t Tuple2[A, B]) Unpack() (A, B) {
	return t.A, t.B
}

// Tuple3 is a group of 3 elements.
type Tuple3[A any, B any, C any] struct {
	A A
	B B
	C C
}

// Unpack returns values contained in tuple.
func (t Tuple3[A, B, C]) Unpack() (A, B, C) {
	return t.A, t.B, t.C
}

// Empty returns an empty value.
func Empty[T any]() T {
	var zero T
	return zero
}

// IsEmpty returns true if argument is a zero value.
func IsEmpty[T comparable](v T) bool {
	var zero T
	return zero == v
}

// IsNotEmpty returns true if argument is not a zero value.
func IsNotEmpty[T comparable](v T) bool {
	var zero T
	return zero != v
}

// ToPtr returns a pointer copy of value.
func ToPtr[T any](x T) *T {
	return &x
}

// FromPtr returns the pointer value or empty.
func FromPtr[T any](x *T) T {
	if x == nil {
		return Empty[T]()
	}

	return *x
}

// ToSlicePtr returns a slice of pointer copy of value.
func ToSlicePtr[T any](collection []T) []*T {
	return slice.Map(collection, func(x T) *T {
		return &x
	})
}

// T2 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T2[A any, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{A: a, B: b}
}

// T2Ptr creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T2Ptr[A any, B any](a A, b B) *Tuple2[A, B] {
	return &Tuple2[A, B]{A: a, B: b}
}

// T3 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T3[A any, B any, C any](a A, b B, c C) Tuple3[A, B, C] {
	return Tuple3[A, B, C]{A: a, B: b, C: c}
}

// T3Ptr creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T3Ptr[A any, B any, C any](a A, b B, c C) *Tuple3[A, B, C] {
	return &Tuple3[A, B, C]{A: a, B: b, C: c}
}

// Zip2 creates a slice of grouped elements, the first of which contains the first elements
// of the given arrays, the second of which contains the second elements of the given arrays, and so on.
// When collections have different size, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip2[A any, B any](a []A, b []B) []Tuple2[A, B] {
	size := slice.Max([]int{len(a), len(b)})

	result := make([]Tuple2[A, B], 0, size)

	for index := 0; index < size; index++ {
		//nolint: errcheck
		_a, _ := slice.Nth(a, index)
		//nolint: errcheck
		_b, _ := slice.Nth(b, index)

		result = append(result, Tuple2[A, B]{
			A: _a,
			B: _b,
		})
	}

	return result
}

// Zip3 creates a slice of grouped elements, the first of which contains the first elements
// of the given arrays, the second of which contains the second elements of the given arrays, and so on.
// When collections have different size, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip3[A any, B any, C any](a []A, b []B, c []C) []Tuple3[A, B, C] {
	size := slice.Max([]int{len(a), len(b), len(c)})

	result := make([]Tuple3[A, B, C], 0, size)

	for index := 0; index < size; index++ {
		//nolint: errcheck
		_a, _ := slice.Nth(a, index)
		//nolint: errcheck
		_b, _ := slice.Nth(b, index)
		//nolint: errcheck
		_c, _ := slice.Nth(c, index)

		result = append(result, Tuple3[A, B, C]{
			A: _a,
			B: _b,
			C: _c,
		})
	}

	return result
}

// Unzip2 accepts an array of grouped elements and creates an array regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip2[A any, B any](tuples []Tuple2[A, B]) ([]A, []B) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)

	for _, tuple := range tuples {
		r1 = append(r1, tuple.A)
		r2 = append(r2, tuple.B)
	}

	return r1, r2
}

// Unzip3 accepts an array of grouped elements and creates an array regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip3[A any, B any, C any](tuples []Tuple3[A, B, C]) ([]A, []B, []C) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)

	for _, tuple := range tuples {
		r1 = append(r1, tuple.A)
		r2 = append(r2, tuple.B)
		r3 = append(r3, tuple.C)
	}

	return r1, r2, r3
}

// AllNil return true when all items are nil or Zero
//
// Deprecated:
// Use function AllEmpty instead.
// This function will be wrong if v is array empty (Ex: AllNil( [ ] int { } ) = false).
func AllNil(v ...any) bool {
	for _, e := range v {
		if e == nil {
			continue
		}
		if !reflect.ValueOf(e).IsZero() {
			return false
		}
	}
	return true
}

func AllEmpty(v ...any) bool {
	for _, e := range v {
		if e == nil {
			continue
		}
		reflectE := reflect.ValueOf(e)
		if reflectE.Kind() == reflect.Slice {
			if reflectE.Len() > 0 {
				return false
			}
			continue
		}
		if !reflect.ValueOf(e).IsZero() {
			return false
		}
	}
	return true
}

// AnyNil return true when at least one item is nil or Zero
func AnyNil(v ...any) bool {
	for _, e := range v {
		if e == nil {
			return true
		}
		if reflect.ValueOf(e).IsZero() {
			return true
		}
	}
	return false
}

func Cast[T any](value any) T {
	var defaultValue T
	if v, ok := value.(T); ok {
		return v
	}
	return defaultValue
}
