package env

import (
	"errors"
)

var (
	ErrNotAField = errors.New("'%s' has no attribute '%s', line %d")
	ErrIndexOutOfRange = errors.New("index out of range, line %d")
	ErrNotArray = errors.New("type %s is not an array, line %d")
	ErrEmptyArray = errors.New("cannot pop from empty array, line %d")
)

// Interface matches all Fizz object structs.
// Type() returns the name of the object
type FizzObject interface {
	Type()  string
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

// Todo: fix reference shit for value and length 
// Stores length value for ease of use. Append elements with += operator.
type Array struct {
	Values *[]interface{}
	Length int
}

func (a Array) Type() string {
	return "array"
}

// Gets value of array at index. Returns error if value is > len(arr) or
// index is less than 0.
func (a Array) Get(index int) (value interface{}, err error) {
	if index >= len(*a.Values) || index < 0 {
		return value, ErrIndexOutOfRange
	}

	return (*a.Values)[index], err
}

// Sets value at given index.
func (a Array) Set(index int, value interface{}) error {
	if index >= a.Length || index < 0 {
		return ErrIndexOutOfRange
	}

	(*a.Values)[index] = value
	return nil
}

// Pushes new value to end of array
func (a *Array) Push(value interface{}) {
	*a.Values = append(*a.Values, value)
	a.Length++
}

// Removes value at end of array and returns it. Returns error if
// length of array is 0.
func (a *Array) Pop() (value interface{}, err error) {
	if a.Length == 0 {
		return nil, ErrEmptyArray
	}

	popped, _ := a.Get(a.Length-1)
	*a.Values = (*a.Values)[:a.Length-1]
	a.Length--
	return popped, nil
}