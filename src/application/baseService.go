package application

import "fmt"

func updateField[T any](source *T, target *T) {
	if source != nil {
		*target = *source
	}
}

func updateEnumField[T fmt.Stringer](currentValue T, newValue *string, getAllValues func() []T) T {
	if newValue == nil {
		return currentValue
	}

	for _, value := range getAllValues() {
		if value.String() == *newValue {
			return value
		}
	}
	return currentValue
}
