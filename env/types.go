package env

import (
	"errors"
	"fmt"
)

var (
	ErrNotAField = errors.New("'%s' has no attribute '%s', line %d")
)

// Interface matches all Fizz object structs.
// Type() returns the name of the object
// Print() returns a printable value of the object
type FizzObject interface {
	Type()  string
	Print() interface{}
}

// Callable object is a function. The origin is the name of the file it
// was defined in. Error returned from Call() is printed as a Fizz error
// and is not a return value.
type Callable struct {
	Name    string
	Origin  string
	Call    func(...interface{}) (interface{}, error)
	NumArgs int
}

func (c Callable) Type() string {
	return "function"
}

func (c Callable) Print() interface{} {
	return c.Name + "()"
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

func (o Object) Print() interface{} {
	str := "{\n"
	for key, value := range o.Fields {
		str += fmt.Sprintf("\t%s: %v\n", key, value)
	}

	return str + "}"
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

// Stores length value for ease of use. Append elements with += operator.
type Array struct {
	Values []interface{}
	Length int
}

func (a Array) Type() string {
	return "array"
}

func (a Array) Print() interface{} {
	return a.Values
}