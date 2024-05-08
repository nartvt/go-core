package math

import (
	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](x T, y T) T {
	if x > y {
		return y
	}
	return x
}

func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// GetValueNotEmptyOrDefault return defaultValue if value is zeroValue, else return value
func GetValueNotEmptyOrDefault[T constraints.Ordered](value, defaultValue T) T {
	var zeroValue T
	if value == zeroValue {
		return defaultValue
	}
	return value
}

// Sum return sum of array number
func Sum[T constraints.Integer | constraints.Float](slices []T) T {
	var sum T
	for _, v := range slices {
		sum += v
	}
	return sum
}

// TernaryOp represent for e1 ? e2 : e3
func TernaryOp[T any](condition bool, ifOutput T, elseOutput T) T {
	if condition {
		return ifOutput
	}
	return elseOutput
}

// TernaryF is a 1 line if/else statement whose options are functions
// Play: https://go.dev/play/p/AO4VW20JoqM
func TernaryF[T any](condition bool, ifFunc func() T, elseFunc func() T) T {
	if condition {
		return ifFunc()
	}

	return elseFunc()
}

// Coalesce returns the first non-empty arguments. Arguments must be comparable.
func Coalesce[T comparable](v ...T) (result T, ok bool) {
	for _, e := range v {
		if e != result {
			result = e
			ok = true
			return
		}
	}

	return
}
