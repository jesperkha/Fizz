package env

import (
	"errors"
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

	return errors.New("variable '" + name + "' is already defined, line %d")
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

	return errors.New("undefined variable '" + name + "', line %d")
}

// Gets variable value by name. Returns error if not defined
func getVariable(env *Environment, name string) (value interface{}, err error) {
	if val, ok := env.Values[name]; ok {
		return val, err
	}

	if env.Parent != nil {
		return getVariable(env.Parent, name)
	}

	return value, errors.New("undefined variable '" + name + "', line %d")
}

func Declare(name string, value interface{}) (err error) {
	return declareVariable(&CurrentEnv, name, value)
}

func Get(name string) (value interface{}, err error) {
	return getVariable(&CurrentEnv, name)
}

func Assign(name string, newVal interface{}) (err error) {
	return assignVariable(&CurrentEnv, name, newVal)
}

// Goes into new scope
func PushScope() {
	newEnv := Environment{Values: entryMap{}}
	newEnv.Parent = &Environment{Parent: CurrentEnv.Parent, Values: CurrentEnv.Values}
	CurrentEnv = newEnv
}

// Adds custom scope for closures
func AddScope(env Environment) {
	env.Parent = &CurrentEnv
	CurrentEnv = env
}

// Goes back to previous scope
func PopScope() {
	parent := CurrentEnv.Parent
	CurrentEnv = *parent
}
