// Package flatter makes array flatten. [[1,2,[3]],4] -> [1,2,3,4]
package flatter

import (
	"errors"
	"log"
)

// Errors
var (
	ErrWrongType = errors.New("flatter: wrong input value")
)

// Int makes some arbitrarily nested arrays of integers into a flat array of integers.
// [[1,2,[3]],4] -> [1,2,3,4]
func Int(arr interface{}) ([]int, error) {

	switch v := arr.(type) {
	case []int:
		return v, nil
	case []interface{}:
		return arbitrarilyInt(v)
	default:
		return nil, ErrWrongType
	}
}

//
func arbitrarilyInt(arr []interface{}) ([]int, error) {

	// first find len of out array for minimize of expensive allocation operation
	len, err := findIntLen(arr)
	if err != nil {
		log.Printf("find arbitrarily int array len error: %v", err)
		return nil, err
	}

	// allocate appropriate size slice for storing out result
	out := make([]int, 0, len)

	// recursively fill out array
	out, err = arbitrarilyIntRecursively(arr, out)
	if err != nil {
		log.Printf("flatter arbitrarily array recursively error: %v", err)
		return nil, err
	}

	return out, nil
}

func arbitrarilyIntRecursively(arr []interface{}, out []int) ([]int, error) {

	for _, iv := range arr {

		switch v := iv.(type) {
		case int:
			out = append(out, v)
		case []int:
			out = append(out, v...)
		case []interface{}:
			var err error
			out, err = arbitrarilyIntRecursively(v, out)
			if err != nil {
				return nil, err
			}
		default:
			return nil, ErrWrongType
		}
	}
	return out, nil
}

// finds len of int items for next allocation memory for new out array
func findIntLen(arr []interface{}) (int, error) {

	l := 0

	for _, iv := range arr {

		switch v := iv.(type) {
		case int:
			l++
		case []int:
			l += len(v)
		case []interface{}:
			l2, err := findIntLen(v)
			if err != nil {
				return 0, err
			}
			l += l2
		default:
			return 0, ErrWrongType
		}
	}
	return l, nil
}
