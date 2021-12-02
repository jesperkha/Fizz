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
