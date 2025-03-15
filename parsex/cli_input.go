package parsex

import "fmt"

type Input struct {
	values map[string]any
	args   []string
}

func newInput() Input {
	return Input{
		values: map[string]any{},
		args:   []string{},
	}
}

func (input Input) Args() []string {
	return input.args
}

func (input Input) Has(key string) bool {
	_, ok := input.values[key]
	return ok
}

func (input Input) String(key string) string {
	if !input.Has(key) {
		return ""
	}
	return input.values[key].(string)
}

func (input Input) Int(key string) int {
	if !input.Has(key) {
		return 0
	}
	return input.values[key].(int)
}

func (input Input) Float(key string) float64 {
	if !input.Has(key) {
		return 0.0
	}
	return input.values[key].(float64)
}

func (input Input) Bool(key string) bool {
	if !input.Has(key) {
		return false
	}
	return input.values[key].(bool)
}

func (input Input) ListOfStrings(key string) []string {
	if !input.Has(key) {
		return []string{}
	}
	items, err := convertToSlice[string](input.values[key])
	if err != nil {
		return []string{}
	}
	return items
}

func (input Input) ListOfInts(key string) []int {
	if !input.Has(key) {
		return []int{}
	}
	items, err := convertToSlice[int](input.values[key])
	if err != nil {
		return []int{}
	}
	return items
}

func (input Input) ListOfFloats(key string) []float64 {
	if !input.Has(key) {
		return []float64{}
	}
	items, err := convertToSlice[float64](input.values[key])
	if err != nil {
		return []float64{}
	}
	return items
}

func (input Input) ListOfBools(key string) []bool {
	if !input.Has(key) {
		return []bool{}
	}
	items, err := convertToSlice[bool](input.values[key])
	if err != nil {
		return []bool{}
	}
	return items
}

func convertToSlice[T comparable](val any) ([]T, error) {
	// Assert the value is a slice of interface{}
	slice, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("value is not a slice")
	}

	items := make([]T, len(slice))
	for i, v := range slice {
		b, ok := v.(T)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a bool", i)
		}
		items[i] = b
	}
	return items, nil
}
