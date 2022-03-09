package env

// Returns pointer to fizz object instance
func NewObject(name string, elements map[string]interface{}) *Object {
	return &Object{Name: name, NumFields: len(elements), Fields: elements}
}

// Returns pointer to fizz callable instance
func NewFunction(name string, numArgs int, callback CallFunction) *Callable {
	return &Callable{Name: name, NumArgs: numArgs, Call: callback}
}

// Returns pointer to fizz array instance
func NewArray(elements []interface{}) *Array {
	return &Array{Values: elements, Length: len(elements)}
}
