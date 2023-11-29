package utils

import "errors"

func StringLengthInRange(min, max int) (func(string) bool, error) {
	if min > max {
		return nil, errors.New("min cannot be larger than max")
	}
	return func(str string) bool {
		l := len(str)
		return l >= min && l <= max
	}, nil
}