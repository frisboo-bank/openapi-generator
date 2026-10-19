package module

import "fmt"

type DependenciesMap[T any] map[string]T

func (m DependenciesMap[T]) Get(name string) (T, error) {
	val, ok := m[name]
	if !ok {
		var zero T
		return zero, fmt.Errorf("dependency %q of type %T not found", name, zero)
	}
	return val, nil
}
