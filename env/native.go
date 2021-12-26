package env

// Native Fizz functions

var StandardEnvironment = Environment{{
	// Returns length of array or string
	"len": Callable{
		NumArgs: 1,
		Call: func(i ...interface{}) (interface{}, error) {
			if arr, ok := i[0].(*Array); ok {
				return float64(arr.Length), nil
			}

			if str, ok := i[0].(string); ok {
				return float64(len(str)), nil
			}

			return -1, ErrNotArray
		},
	},

	// Append new element to end of array
	"push": Callable{
		NumArgs: 2,
		Call: func(i ...interface{}) (interface{}, error) {
			if arr, ok := i[0].(*Array); ok {
				arr.Push(i[1])
				return nil, nil
			}

			return -1, ErrNotArray
		},
	},

	// Remove last element and return it
	"pop": Callable{
		NumArgs: 1,
		Call: func(i ...interface{}) (interface{}, error) {
			if arr, ok := i[0].(*Array); ok {
				return arr.Pop()
			}

			return -1, ErrNotArray
		},
	},
}}
