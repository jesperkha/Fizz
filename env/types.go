package env

import (
	"errors"
)

var (
	ErrNotAField       = errors.New("'%s' has no attribute '%s', line %d")
	ErrIndexOutOfRange = errors.New("index out of range, line %d")
	ErrNotArray        = errors.New("type %s is not an array, line %d")
	ErrEmptyArray      = errors.New("cannot pop from empty array, line %d")
)

// Performs recursive equality check for objects and arrays.
// Other cases returns standard equality check.
func Equal(left, right interface{}) bool {
	l, r := "", ""
	if lo, ok := left.(FizzObject); ok {
		l = lo.Type()
	}

	if ro, ok := right.(FizzObject); ok {
		r = ro.Type()
	}

	if l == "array" && r == "array" {
		a, _ := left.(*Array)
		b, _ := right.(*Array)
		return a.IsEqual(b)
	}

	if l == "object" && r == "object" {
		a, _ := left.(*Object)
		b, _ := right.(*Object)
		return a.IsEqual(b)
	}

	return left == right
}

// Interface matches all Fizz object structs.
// Type() returns the name of the object
type FizzObject interface {
	Type() string
}

type CallFunction func(...interface{}) (interface{}, error)

// Callable object is a function. The origin is the name of the file it
// was defined in. Error returned from Call() is printed as a Fizz error
// and is not a return value.
type Callable struct {
	Name    string
	Origin  string
	Call    CallFunction
	NumArgs int
}

func (c *Callable) Type() string {
	return "function"
}

// Object with n fields. Name is the name of the constructor, not the
// instance. File imports are also objects.
type Object struct {
	Fields    map[string]interface{}
	NumFields int
	Name      string
}

func (o *Object) Type() string {
	return "object"
}

// Equality check for objects
func (o *Object) IsEqual(a *Object) bool {
	if o.NumFields != a.NumFields {
		return false
	}

	for k, v := range o.Fields {
		ov, err := a.Get(k)
		if err != nil {
			return false
		}

		if !Equal(v, ov) {
			return false
		}
	}

	return true
}

// Gets value from object. Used for getter syntax "name.value"
func (o *Object) Get(name string) (value interface{}, err error) {
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

// Compare two arrays
func (a *Array) IsEqual(o *Array) bool {
	if a.Length != o.Length {
		return false
	}

	for i, n := range a.Values {
		if !Equal(o.Values[i], n) {
			return false
		}
	}

	return true
}

// Gets value of array at index. Returns error if value is > len(arr) or
// index is less than 0.
func (a Array) Get(index int) (value interface{}, err error) {
	if index >= len(a.Values) || index < 0 {
		return value, ErrIndexOutOfRange
	}

	return a.Values[index], err
}

// Sets value at given index.
func (a Array) Set(index int, value interface{}) error {
	if index >= a.Length || index < 0 {
		return ErrIndexOutOfRange
	}

	a.Values[index] = value
	return nil
}

// Pushes new value to end of array
func (a *Array) Push(value interface{}) {
	a.Values = append(a.Values, value)
	a.Length++
}

// Removes value at end of array and returns it. Returns error if
// length of array is 0.
func (a *Array) Pop() (value interface{}, err error) {
	if a.Length == 0 {
		return nil, ErrEmptyArray
	}

	popped, _ := a.Get(a.Length - 1)
	a.Values = a.Values[:a.Length-1]
	a.Length--
	return popped, nil
}
