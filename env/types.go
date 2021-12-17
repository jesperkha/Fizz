package env

import (
	"errors"
)

var (
	ErrNotAField = errors.New("'%s' has no attribute '%s', line %d")
)

// Interface matches all Fizz object structs.
// Type() returns the name of the object
type FizzObject interface {
	Type() string
}

// Callable object is a function. The origin is the name of the file it
// was defined in. Error returned from Call() is printed as a Fizz error
// and is not a return value.
type Callable struct {
	Origin  string
	Call    func(...interface{}) (interface{}, error)
	NumArgs int
}

func (c Callable) Type() string {
	return "function"
}

// Object with n fields. Name is the name of the constructor, not the
// instance. File imports are also objects.
type Object struct {
	Fields    map[string]interface{}
	NumFields int
	Name      string
}

func (o Object) Type() string {
	return "object"
}

// Gets value from object. Used for getter syntax "name.value"
func (o Object) Get(name string) (value interface{}, err error) {
	if val, ok := o.Fields[name]; ok {
		return val, err
	}

	return value, ErrNotAField
}

// Reassigns value to object. Does not declare since object have a
// constant number of fields. Used for setter syntax "name.value = n"
func (o *Object) Set(name string, value interface{}) (err error) {
	if _, ok := o.Fields[name]; ok {
		o.Fields[name] = value
		return err
	}

	return ErrNotAField
}