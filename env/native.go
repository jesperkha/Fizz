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

type Object struct {
	Fields 	  map[string]interface{}
	NumFields int
}

func (o Object) Type() string {
	return "object"
}