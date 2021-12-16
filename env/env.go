package env

import (
	"errors"
)

var (
	ThrowEnvironment = true
)

type valueMap map[string]interface{}

// List of 'scopes'. Index 0 is always the current scope and when looping the
// order will be from low to high level scopes.
type Environment []valueMap

var currentEnv = Environment{{}}
var tempEnv = currentEnv

// Creates new environment, replacing the old one. Returns old environment. For
// testing, its not necessary to get rid of the old env, hence the option to not
// remove it.
func NewEnvironment() Environment {
	if !ThrowEnvironment {
		return currentEnv
	}

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
// scopes is possible. Returns error if name is not defined anywhere.
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

// Adds environment from imported file to current env. Values are passed as an object.
// The fields are the values in the global scope of the environment (index 0).
func AddImportedFile(name string, env Environment) error {
	return Declare(name, Object{
		Name: name,
		NumFields: len(env[0]),
		Fields: env[0],
		File: true,
	})
}

// Gets a snapshot of the current env. Used for composing closures for functions and
// also file imports. Unsafe: will get a temp env if there is one active.
func GetCurrentEnv() Environment {
	return []valueMap(currentEnv)
}

// Sets a new temporary envirnoment. Used for closures since envs are not passed as
// arguments to any functions in this file. Is discarded upon calling PopTempEnv().
func PushTempEnv(env Environment) {
	tempEnv = currentEnv
	currentEnv = env
}

// Unsafe: does not check if there is a current temp env or not, however, its use is
// hardcoded and will not be called when there is no temporary environment.
func PopTempEnv() {
	currentEnv = tempEnv
}