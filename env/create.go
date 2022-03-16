package env

func NewObject(name string, elements map[string]interface{}) *Object {
	return &Object{Name: name, NumFields: len(elements), Fields: elements}
}

func NewFunction(name string, numArgs int, callback CallFunction) *Callable {
	return &Callable{Name: name, NumArgs: numArgs, Call: callback}
}

func NewArray(elements []interface{}) *Array {
	return &Array{Values: elements, Length: len(elements)}
}
