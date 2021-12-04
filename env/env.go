package env

import (
	"errors"
)

// Todo: redesign entire env model. every running instance of the interpreter should be 100% self contained
type entryMap map[string]interface{}

type Environment struct {
	Values entryMap
	Parent *Environment
	Name string
}

var currentEnv = Environment{Values: entryMap{}}

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
	return declareVariable(&currentEnv, name, value)
}

func Get(name string) (value interface{}, err error) {
	return getVariable(&currentEnv, name)
}

func Assign(name string, newVal interface{}) (err error) {
	return assignVariable(&currentEnv, name, newVal)
}

// Goes into new scope
func PushScope() {
	newEnv := Environment{Values: entryMap{}}
	newEnv.Parent = &Environment{Parent: currentEnv.Parent, Values: currentEnv.Values}
	currentEnv = newEnv
}

// Adds custom scope for closures
func AddScope(env Environment) {
	env.Parent = &currentEnv
	currentEnv = env
}

// Goes back to previous scope
func PopScope() {
	parent := currentEnv.Parent
	currentEnv = *parent
}
