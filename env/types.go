package env

import (
	"errors"
)

var (
	ErrNotAField = errors.New("'%s' has no attribute '%s', line %d")
)

type FizzObject interface {
	Type() string
}

type Callable struct {
	Origin  string // Name of file function was defined in for error tracing
	Call    func(...interface{}) (interface{}, error)
	NumArgs int
}

func (c Callable) Type() string {
	return "function"
}

type Object struct {
	Fields    map[string]interface{}
	NumFields int
	Name      string
	File 	  bool
}

func (o Object) Type() string {
	return "object"
}

func (o Object) Get(name string) (value interface{}, err error) {
	if val, ok := o.Fields[name]; ok {
		return val, err
	}

	return value, ErrNotAField
}

func (o *Object) Set(name string, value interface{}) (err error) {
	if _, ok := o.Fields[name]; ok {
		o.Fields[name] = value
		return err
	}

	return ErrNotAField
}