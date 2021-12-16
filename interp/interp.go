package interp

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jesperkha/Fizz/env"
	"github.com/jesperkha/Fizz/lexer"
	"github.com/jesperkha/Fizz/stmt"
	"github.com/jesperkha/Fizz/util"
)

var (
	ErrFileNotFound = errors.New("cannot find file with name: '%s'")
	ErrNonFizzFile  = errors.New("cannot run non-Fizz file")
	ErrCircularImport = errors.New("circular import not allowed, %s <-> %s")
)

// Stores import pairs to check for import cycles. Duplicate erntries indicate
// a circular import and an error is raised.
var importPairs = util.UniquePairs{}

// Interperates string of code. The string is tokenized in the lexer package
// where each token has a type, line, lexeme, and literal value. The parsed
// tokens are passed to the statement parser which looks for statements
// (identified by trailing semicolon or curly braces).

// The statements are parsed as statement tokens with a type, name, and
// expression values for the expression associated with the statement. For
// example: a variable declaration statement has an expression property for
// the initial value that is found after the equal sign.

// The statements are then executed and, when doing so, statement expressions
// are evaluated. Variable values are also assigned and manipulated in the
// variable environment found in the env package.

func Interperate(filename string, input string) (e env.Environment, err error) {
	// Parses input characters into lexical tokens for single and double symbols,
	// identifiers, and keywords.
	lexicalTokens, err := lexer.GetTokens(input)
	if err != nil {
		return e, err
	}

	// Lexical tokens are analysed and put into statement tokens. These statements
	// contain all the information they need for execution and error handling.
	statements, err := stmt.ParseStatements(lexicalTokens)
	if err != nil {
		return e, err
	}

	// File imports are handled after parsing the statements, not when executing them.
	// This means that all imports are run before anything else; they are "hoisted".
	// Even declaring a variable with the same name as the file before importing it will
	// raise an error as the file is imported before the variable is created.
	for _, s := range statements {
		if s.Type != stmt.Import {
			continue
		}
		
		// Checks for circular imports. Add() returns true if the pair already exists.
		name := util.GetPlainFilename(s.Name)
		if this := util.GetPlainFilename(filename); importPairs.Add(name, this) {
			return e, fmt.Errorf(ErrCircularImport.Error(), name, this)
		}

		e, err = RunFile(s.Name + ".fizz")
		if err != nil {
			return e, err
		}
		
		// Adds the global environment of the imported file to the env of the current one.
		// It is added as an object instance with the name of the file without the fizz suffix.
		if err = env.AddImportedFile(name, e); err != nil {
			return e, err
		}
	}
	
	// Set origin point for function declarations. This makes sure that errors give
	// the correct filename when printed.
	stmt.CurrentOrigin = filename
	
	// Finally executes statement tokens. This is the only step that has any effect
	// on the actual input program as the others were just breaking it up into usable
	// pieces. While the interpreter is still running, the values of variables will be
	// remembered as the environments are never reset at runtime.
	err = stmt.ExecuteStatements(statements)
	return env.NewEnvironment(),  err
}

// Runs a fizz file. Imports are run as files and the environment is extracted and
// packaged into a namespace. Said namespace is put into the file it was imported from,
// which means if "main.fizz" imports "other.fizz", the main file also imports all
// of the files imported in "other.fizz".
func RunFile(filename string) (e env.Environment, err error) {
	if !strings.Contains(filename, ".") {
		filename = filename + ".fizz"
	}

	if !strings.HasSuffix(filename, ".fizz") {
		return e, ErrNonFizzFile
	}

	if byt, err := os.ReadFile(filename); err == nil {
		e, err = Interperate(filename, string(byt))
		return e, util.WrapFilename(filename, err)
	}

	// Unsafe: assumes path error
	return e, fmt.Errorf(ErrFileNotFound.Error(), filename)
}
