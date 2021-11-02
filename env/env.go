package env

import (
	"errors"
	"fmt"
)

var (
	ErrUndefinedVariable = errors.New("undefined variable '%s', line %d")
	ErrAlreadyDefined	 = errors.New("variable '%s' is already defined, line %d")
)

type entryMap map[string]interface{}

type Environment struct {
	Values entryMap
	Parent *Environment
}

var CurrentEnv = Environment{Values: entryMap{}}

// Adds a new key value pair for specified variable name. Returns error if already defined
func declareVariable(env *Environment, name string, value interface{}) (err error) {
	if _, ok := env.Values[name]; !ok {
		env.Values[name] = value
		return err
	}

	return fmt.Errorf(ErrAlreadyDefined.Error(), name)
}

// Reassigns variable value if exists. Returns error otherwise
func assignVariable(env *Environment, name string, newVal interface{}) (err error) {
	if _, ok := env.Values[name]; ok {
		env.Values[name] = newVal
		return err
	}

	if env.Parent != nil {
		return assignVariable(env.Parent, name, newVal)
	}

	return fmt.Errorf(ErrUndefinedVariable.Error(), name)
}

// Gets variable value by name. Returns error if not defined
func getVariable(env *Environment, name string) (value interface{}, err error) {
	if val, ok := env.Values[name]; ok {
		return val, err
	}

	if env.Parent != nil {
		return getVariable(env.Parent, name)
	}

	return value, fmt.Errorf(ErrUndefinedVariable.Error(), name)
}

// Declares to current scope enviromnent table
func Declare(name string, value interface{}) (err error) {
	return declareVariable(&CurrentEnv, name, value)
}

// Gets value from current scope environment table
func Get(name string) (value interface{}, err error) {
	return getVariable(&CurrentEnv, name)
}

// Assigns new value to variable in current scope
func Assign(name string, newVal interface{}) (err error) {
	return assignVariable(&CurrentEnv, name, newVal)
}

// Goes into new scope
func PushScope() {
	newEnv := Environment{Values: entryMap{}}
	newEnv.Parent = &Environment{Parent: CurrentEnv.Parent, Values: CurrentEnv.Values}
	CurrentEnv = newEnv
}

// Goes back to previous scope
func PopScope() {
	parent := CurrentEnv.Parent
	CurrentEnv = *parent
}