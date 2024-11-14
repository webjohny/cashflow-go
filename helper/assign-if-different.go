package helper

import "reflect"

func AssignIfDifferent[T comparable](target *T, newValue T) {
	if !reflect.DeepEqual(*target, newValue) {
		*target = newValue
	}
}
