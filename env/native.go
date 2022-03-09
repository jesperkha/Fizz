package env

var StandardEnvironment = Environment{{
	// Returns length of array or string
	"len": NewFunction("len", 1, func(i ...interface{}) (interface{}, error) {
		if arr, ok := i[0].(*Array); ok {
			return float64(arr.Length), nil
		}

		if str, ok := i[0].(string); ok {
			return float64(len(str)), nil
		}

		return -1, ErrNotArray
	}),

	// Append new element to end of array
	"push": NewFunction("push", 2, func(i ...interface{}) (interface{}, error) {
		if arr, ok := i[0].(*Array); ok {
			arr.Push(i[1])
			return nil, nil
		}

		return -1, ErrNotArray
	}),

	// Remove last element and return it
	"pop": NewFunction("pop", 1, func(i ...interface{}) (interface{}, error) {
		if arr, ok := i[0].(*Array); ok {
			return arr.Pop()
		}

		return -1, ErrNotArray
	}),
}}
