package lib

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/util"
)

type i interface{}
type FuncMap map[string]interface{}

var (
	LibList   = map[string]FuncMap{}
	ErrNotLib = errors.New("%s is not a library")
)

// Adds inclusion map to global include list if lib name is required
func Add(libName string, functions FuncMap) {
	if _, ok := LibList[libName]; ok {
		util.ErrorAndExit(fmt.Errorf("duplicate library name '%s'", libName))
	}

	LibList[libName] = functions
}

// Imports all included libs to running fizz process
func IncludeLibraries(includes []string) error {
	for _, name := range includes {
		// Check if name is in list of valid libs
		library, ok := LibList[name]
		if !ok {
			return fmt.Errorf(ErrNotLib.Error(), name)
		}

		// Declare all functions in library. Push new scope, declare functions
		// and then get env before popping scope again. This simulates a file
		// import and the functions are added the same way as normal imports.
		env.PushScope()
		for funcName, f := range library {
			// Cache values because they are used in the Call member, which will use
			// the closured version of the variables, and they will always be the last
			// used in the loop. Cache makes sure its always the ones at definition.
			nameCache, funcCache := funcName, f

			// Create function. -1 ignores number of args in parsing
			callable := env.Callable{
				NumArgs: -1,
				Origin:  name,
				Name:    funcName,
				Call: func(args ...interface{}) (interface{}, error) {
					return CallFunc(nameCache, funcCache, args)
				},
			}

			// Declare to scope it was required in
			if err := env.Declare(funcName, callable); err != nil {
				return err
			}
		}

		// Get env and pop. Add as import
		defEnv := env.GetCurrentEnv()
		env.PopScope()
		if err := env.AddImportedFile(name, defEnv); err != nil {
			return err
		}
	}

	return nil
}

// Calls library function from fizz process. Uses reflect to get arg number and
// types from function. This means the types of params matter (but only here).
// Returns error if args or types dont match. Name of function and name given do
// not need to match.
func CallFunc(name string, function i, args []interface{}) (val i, err error) {
	// Get function as type
	f := reflect.ValueOf(function)

	// Check if num args are valid
	numArgs := f.Type().NumIn()
	gotArgs := len(args)
	if numArgs != gotArgs {
		s := fmt.Sprintf("%s() expected %d args, got %d", name, numArgs, gotArgs)
		return val, errors.New(s + ", line %d")
	}

	// Convert args to reflect.Value
	argsIn := make([]reflect.Value, numArgs)
	for idx, value := range args {
		// Arg value types must match for lib functions
		paramType := f.Type().In(idx)
		argType := reflect.TypeOf(value)
		// Interface param type doesnt need type check.
		// Unsafe: i might not be defined as interface
		if paramType != argType && paramType.Name() != "i" {
			// Get fizz names for types
			expect := util.GetLibType(paramType.Name())
			got := util.GetType(value)

			// Do crazy error shenanigans to avoid writing twice
			s := "%s() expected arg %d to be %s, got %s"
			e := fmt.Sprintf(s, name, idx+1, expect, got)
			return val, errors.New(e + ", line %d")
		}

		argsIn[idx] = reflect.ValueOf(value)
	}

	// Call function with args
	res := f.Call(argsIn)
	value := res[0].Interface()
	// Get error. Return interface value is techically never nil so an
	// explicit check is performed.
	if res[1].IsNil() {
		err = nil
	} else {
		err = res[1].Interface().(error)
	}

	// Returned value must be error based in map type
	return value, err
}
