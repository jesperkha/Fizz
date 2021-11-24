package env

import "log"

type Callable struct {
	Call    func(...interface{}) (interface{}, error)
	NumArgs int
}

func init() {
	Declare("log", Callable{
		NumArgs: 1,
		Call: func(i ...interface{}) (interface{}, error) {
			log.Println(i...)
			return 1.0, nil
		},
	})
}
