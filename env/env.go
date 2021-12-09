package env

import (
	"errors"
)

type valueMap map[string]interface{}

// List of 'scopes'. Index 0 is always the current scope and when looping the
// order will be from low to high level scopes.
type Environment []valueMap

var currentEnv = Environment{{}}

// Creates new environment, replacing the old one. Returns old environment.
func NewEnvironment() Environment {
	oldEnv := currentEnv
	currentEnv = Environment{{}}
	return oldEnv
}

// Declares value to name in current scope. This allows for overriding global
// variable names for local scopes. Returns error if name is already declared.
func Declare(name string, value interface{}) error {
	curScope := currentEnv[0]
	if _, ok := curScope[name]; !ok {
		curScope[name] = value
		return nil
	}

	return errors.New("variable '" + name + "' is already defined, line %d")
}

// Assigns value to name. If name is not defined in current scope the parent
// scopes are checked. Therefore, reassignment of global variables in local
// scopes i possible. Returns error if name is not defined anywhere.
func Assign(name string, value interface{}) error {
	for _, scope := range currentEnv {
		if _, ok := scope[name]; ok {
			scope[name] = value
			return nil
		}
	}

	return errors.New("undefined variable '" + name + "', line %d")
}

// Gets the value assigned to name. If the name is not defined in the current
// scope the parent scopes are checked. Gets the first instance of name. Returns
// error if name is not defined anywhere.
func Get(name string) (value interface{}, err error) {
	for _, scope := range currentEnv {
		if value, ok := scope[name]; ok {
			return value, nil
		}
	}

	return value, errors.New("undefined variable '" + name + "', line %d")
}

// Puts new scope at beginning of slice, effectivly setting the previous scope
// as the parent of the new. (slice is reverse stack)
func PushScope() {
	currentEnv = append([]valueMap{{}}, currentEnv...)
}

// Removes first element in slice, meaning the parent scope is set to the current.
// Unsafe: if the length of the slice is 1, pop will panic. However, the use of Push()
// and Pop() is hardcoded and will never cause a pop of a scope list smaller than 2.
func PopScope() {
	if len(currentEnv) < 2 {
		panic("env: popped scope list of length < 2")
	}

	currentEnv = currentEnv[1:]
}
