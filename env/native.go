package env

// Native Fizz functions

// Todo: push, pop
var StandardEnvironment = Environment{
	{
		// Gets length of array
		"len": Callable {
		NumArgs: 1,
		Call: func(i ...interface{}) (interface{}, error) {
			if arr, ok := i[0].(Array); ok {
				return float64(arr.Length), nil
			}

			return -1, ErrNotArray
		},
		},
	},
}