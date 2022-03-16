package env

var StandardEnvironment = Environment{{
	"len": NewFunction("len", 1, func(i ...interface{}) (interface{}, error) {
		if arr, ok := i[0].(*Array); ok {
			return float64(arr.Length), nil
		}

		// Strings too
		if str, ok := i[0].(string); ok {
			return float64(len(str)), nil
		}

		return -1, ErrNotArray
	}),

	"push": NewFunction("push", 2, func(i ...interface{}) (interface{}, error) {
		if arr, ok := i[0].(*Array); ok {
			arr.Push(i[1])
			return nil, nil
		}

		return -1, ErrNotArray
	}),

	"pop": NewFunction("pop", 1, func(i ...interface{}) (interface{}, error) {
		if arr, ok := i[0].(*Array); ok {
			return arr.Pop()
		}

		return -1, ErrNotArray
	}),
}}
