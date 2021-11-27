package env

type FizzObject interface {
	Type()  string
}

type Callable struct {
	Call    func(...interface{}) (interface{}, error)
	NumArgs int
}

func (c Callable) Type() string {
	return "function"
}

func init() {
	Declare("TEST", Callable{
		NumArgs: 0,
		Call: func(i ...interface{}) (interface{}, error) {
			return 1.0, nil
		},
	})
}
