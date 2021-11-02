package env

import (
	"errors"
	"fmt"
)

var (
	ErrUndefinedVariable = errors.New("undefined variable '%s', line %d")
	ErrAlreadyDefined	 = errors.New("variable '%s' is already defined, line %d")
)

type Environment map[string]interface{}
var GlobalEnv = Environment{}
var CurrentScope = GlobalEnv

// Adds a new key value pair for specified variable name. Returns error if already defined
func declareVariable(env Environment, name string, value interface{}) (err error) {
	if _, ok := env[name]; !ok {
		env[name] = value
		return err
	}

	return fmt.Errorf(ErrAlreadyDefined.Error(), name)
}

// Gets variable value by name. Returns error if not defined
func getVariable(env Environment, name string) (value interface{}, err error) {
	if val, ok := env[name]; ok {
		return val, err
	}

	return value, fmt.Errorf(ErrUndefinedVariable.Error(), name)
}

// Declares to global enviromnent table
func Declare(name string, value interface{}) (err error) {
	return declareVariable(GlobalEnv, name, value)
}

// Gets value from global environment table
func Get(name string) (value interface{}, err error) {
	return getVariable(GlobalEnv, name)
}